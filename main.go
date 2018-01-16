package main

import (
	"fmt"
	"log"
	"sort"

	env "github.com/segmentio/go-env"
	dropbox "github.com/tj/go-dropbox"
)

const PERPAGE = 1000
const SHOWRESULTS = 20

type sortBySize []*dropbox.SearchMatch

func (a sortBySize) Len() int           { return len(a) }
func (a sortBySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortBySize) Less(i, j int) bool { return a[i].Metadata.Size > a[j].Metadata.Size }

type sortByTime []*dropbox.SearchMatch

func (a sortByTime) Len() int      { return len(a) }
func (a sortByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortByTime) Less(i, j int) bool {
	return a[i].Metadata.ClientModified.Format("2006-01-02") > a[j].Metadata.ClientModified.Format("2006-01-02")
}

func main() {
	c := client()

	var files []*dropbox.SearchMatch
	var i, start uint64 = 0, 0

	for {
		start = i * PERPAGE

		out, err := c.Files.Search(&dropbox.SearchInput{
			Path:       "/",
			Query:      "_201*_",
			MaxResults: PERPAGE,
			Start:      start,
		})

		if err != nil {
			panic(err)
		}

		log.Printf("Page: %d / Results: %d", i+1, len(out.Matches))

		for _, file := range out.Matches {
			files = append(files, file)
		}

		if !out.More {
			break
		}

		i++
	}

	if len(files) == 0 {
		return
	}

	fmt.Println("\nSorted by size DESC")
	sort.Sort(sortBySize(files))

	i = 0
	for _, file := range files {
		i++

		size := float64(file.Metadata.Size) / 1024 / 1024

		if i > SHOWRESULTS-1 {
			break
		}
		fmt.Printf("File: %s = %f MB\n", file.Metadata.PathDisplay, size)
	}

	fmt.Println("\nSorted by date DESC")
	sort.Sort(sortByTime(files))

	i = 0
	for _, file := range files {
		i++

		size := float64(file.Metadata.Size) / 1024 / 1024

		if i > SHOWRESULTS-1 {
			break
		}

		fmt.Printf("File (created at: %s): %s = %f MB \n", file.Metadata.ClientModified.Format("2006-01-02"), file.Metadata.PathDisplay, size)
	}
}

func client() *dropbox.Client {
	token, err := env.Get("DROPBOX_ACCESS_TOKEN")

	if err != nil {
		panic("No token provided. Run `export DROPBOX_ACCESS_TOKEN=arUfR.......Cc`")
	}

	return dropbox.New(dropbox.NewConfig(token))
}
