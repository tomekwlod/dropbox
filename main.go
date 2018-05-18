package dropbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	dropbox "github.com/tj/go-dropbox"
	config "github.com/tomekwlod/dropbox/config"
)

const perpage = 1000
const uploadSizeLimit = 140

var uploadSessionToken string

type FilesResult struct {
	Size     uint64    `json:"size"`
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
}
type StorageResult struct {
	Used                uint64  `json:"used"`
	AllocationUsed      uint64  `json:"allocation_used"`
	AllocationAllocated uint64  `json:"allocation_allocated"`
	PercentUsed         float32 `json:"percent_used"`
}
type UploadResult struct {
	ContentHash string `json:"content_hash"`
	Rev         string `json:"rev"`
	Size        int64  `json:"size"`
	Name        bool   `json:"name"`
	PathLower   string `json:"path_lower"`
	PathDisplay string `json:"path_display"`
	ID          string `json:"id"`
}
type UploadInput struct {
	Path           string `json:"path"`
	LocalFile      string `json:"-"`
	Mode           string `json:"mode"`
	AutoRename     bool   `json:"autorename"`
	Mute           bool   `json:"mute"`
	ClientModified string `json:"client_modified,omitempty"`
}
type Session struct {
	Id string `json:"session_id"`
}
type SessionAppend struct {
	Cursor Cursor `json:"cursor"`
	Close  bool   `json:"close"`
}
type SessionFinish struct {
	Cursor Cursor      `json:"cursor"`
	Commit UploadInput `json:"commit"`
}
type Cursor struct {
	SessionID string `json:"session_id"`
	Offset    int64  `json:"offset"`
}

func Search(term string) []*dropbox.SearchMatch {
	c := config.DropboxClient()

	var files []*dropbox.SearchMatch
	var i, start uint64 = 0, 0

	for {
		start = i * perpage

		out, err := c.Files.Search(&dropbox.SearchInput{
			Path:       "/",
			Query:      term,
			MaxResults: perpage,
			Start:      start,
		})

		if err != nil {
			panic(err)
		}

		for _, file := range out.Matches {
			files = append(files, file)
		}

		if !out.More {
			break
		}

		i++
	}

	return files
}

func AllFiles() []*FilesResult {
	c := config.DropboxClient()

	var files []*FilesResult

	in := dropbox.ListFolderInput{Path: "", Recursive: true}
	list, err := c.Files.ListFolder(&in)

	if err != nil {
		panic(err)
	}

	for _, v := range list.Entries {
		if v.Size > 0 {
			files = append(files, &FilesResult{Size: v.Size, Name: v.PathDisplay, Modified: v.ClientModified})
		}
	}

	if !list.HasMore {
		return files
	}

	cursor := list.Cursor
	for {
		inc := dropbox.ListFolderContinueInput{Cursor: cursor}
		list, err := c.Files.ListFolderContinue(&inc)

		if err != nil {
			panic(err)
		}

		if !list.HasMore {
			break
		}

		for _, v := range list.Entries {
			if v.Size > 0 {
				files = append(files, &FilesResult{Size: v.Size, Name: v.PathDisplay, Modified: v.ClientModified})
			}
		}

		cursor = list.Cursor
	}

	return files
}

func Delete(term string) bool {
	c := config.DropboxClient()

	_, err := c.Files.Delete(&dropbox.DeleteInput{
		Path: term,
	})

	if err != nil {
		return false
	}

	return true
}

func Storage() *StorageResult {
	c := config.DropboxClient()

	out, err := c.Users.GetSpaceUsage()

	if err != nil {
		return nil
	}

	result := &StorageResult{
		Used:                out.Used,
		AllocationUsed:      out.Allocation.Used,
		AllocationAllocated: out.Allocation.Allocated,
		PercentUsed:         float32((100 * float32(out.Used)) / float32(out.Allocation.Allocated)),
	}

	return result
}

func Download(filename, target string) (err error) {
	c := config.DropboxClient()

	body, err := c.Files.Download(&dropbox.DownloadInput{
		Path: filename,
	})
	if err != nil {
		return err
	}

	// mkdir -p location/to/create
	os.MkdirAll(strings.Replace(target, filepath.Base(target), "", -1), os.ModePerm)

	// save body to file
	outFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, body.Body)
	if err != nil {
		panic(err)
	}

	return err
}

// needs to also return an error (not only here but in all of the functions)
func Upload(in *UploadInput) (result *UploadResult) {
	f, err := os.Open(in.LocalFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fileInfo, _ := f.Stat()
	fileSize := fileInfo.Size()

	if fileSize/1024/1024 < uploadSizeLimit {
		return UploadSmall(in)
	}

	// calculate total number of parts the file will be chunked into
	const fileChunk = (uploadSizeLimit * (1 << 20)) // 1 * (1 << 20) = 1 MB, 140 * (1 << 20) = 140MB
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	fmt.Printf("Splitting file into %d pieces.\n", totalPartsNum)

	sentSoFar := 0
	for i := uint64(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(float64(fileChunk), float64(fileSize-int64(int(i)*fileChunk))))
		chunk := make([]byte, partSize)

		// very importan
		f.Read(chunk)

		switch i {
		case 0:
			log.Println("FIRST CALL")
			// first call opening the dropbox session
			body := fileTransfer("https://content.dropboxapi.com/2/files/upload_session/start", bytes.NewReader(chunk), map[string]bool{"close": false})
			var dat map[string]string
			if err := json.Unmarshal(body.Bytes(), &dat); err != nil {
				panic(err)
			}
			uploadSessionToken = dat["session_id"]

			break
		case totalPartsNum - 1:
			log.Println("FINAL CALL")
			// finising the transfer
			params := &SessionFinish{}
			params.Commit = *in
			params.Cursor.SessionID = uploadSessionToken
			params.Cursor.Offset = int64(sentSoFar)

			body := fileTransfer("https://content.dropboxapi.com/2/files/upload_session/finish", bytes.NewReader(chunk), params)
			// var resp map[string]interface{}
			// if err := json.Unmarshal(body.Bytes(), &resp); err != nil {
			// 	panic(err)
			// }
			// fmt.Printf("RESPONSE: %+v", resp)

			json.Unmarshal(body.Bytes(), &result)

			break
		default:
			log.Printf("MIDDLE CALL (%s)\n", uploadSessionToken)
			// middle chunks transfer
			params := &SessionAppend{}
			params.Close = false
			params.Cursor.SessionID = uploadSessionToken
			params.Cursor.Offset = int64(sentSoFar)

			fileTransfer("https://content.dropboxapi.com/2/files/upload_session/append_v2", bytes.NewReader(chunk), params)

			break
		}

		sentSoFar = sentSoFar + partSize
	}

	return
}

func UploadSmall(in *UploadInput) (result *UploadResult) {
	f, err := os.Open(in.LocalFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	body := fileTransfer("https://content.dropboxapi.com/2/files/upload", f, in)
	json.Unmarshal(body.Bytes(), &result)

	return
}

func fileTransfer(uri string, chunk io.Reader, params interface{}) (body *bytes.Buffer) {
	dbxToken, err := config.DropboxToken()
	if err != nil {
		panic("No token provided. Run `export DROPBOX_ACCESS_TOKEN=arUfR.......Cc`")
	}

	req, err := http.NewRequest("POST", uri, chunk)
	req.Header.Add("Authorization", "Bearer "+dbxToken)
	req.Header.Add("Content-Type", "application/octet-stream")

	p, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Dropbox-API-Arg", string(p))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body = &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
		panic(err)
	}

	return
}
