package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	dropbox "github.com/tj/go-dropbox"
	"github.com/tomekwlod/dropboxCleaner/files"
)

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
	term := "_201*_"

	if len(os.Args) > 1 {
		term = os.Args[1]
	}

	log.Printf("Using term: `%s`\n", term)

	files := files.Search(term)

	if len(files) == 0 {
		log.Printf("No files found for a term: `%s`\n", term)

		return
	}

	fmt.Println("\nSorted by size DESC")
	sort.Sort(sortBySize(files))

	i := 0
	for _, file := range files {
		i++

		size := float64(file.Metadata.Size) / 1024 / 1024

		if i > SHOWRESULTS-1 {
			break
		}

		fmt.Printf("File (created at: %s): %s = %f MB \n", file.Metadata.ClientModified.Format("2006-01-02"), file.Metadata.PathDisplay, size)
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
