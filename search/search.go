package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	godropbox "github.com/tj/go-dropbox"
	"github.com/tomekwlod/dropbox"
)

const SHOWRESULTS = 20

type sortBySize []*godropbox.SearchMatch

func (a sortBySize) Len() int           { return len(a) }
func (a sortBySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortBySize) Less(i, j int) bool { return a[i].Metadata.Size > a[j].Metadata.Size }

type sortByTime []*godropbox.SearchMatch

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

	log.Printf("Searching using term: `%s`\n", term)

	files := dropbox.Search(term)

	if len(files) == 0 {
		log.Printf("No files found for a term: `%s`\n", term)

		return
	}

	log.Printf("%d files found\n", len(files))

	fmt.Println("\nSorted by size DESC")
	sort.Sort(sortBySize(files))

	i := 1
	for _, file := range files {
		size := float64(file.Metadata.Size) / 1024 / 1024

		fmt.Printf("File (created at: %s): %s = %f MB \n", file.Metadata.ClientModified.Format("2006-01-02"), file.Metadata.PathDisplay, size)

		if i >= SHOWRESULTS {
			break
		}

		i++
	}

	fmt.Println("\nSorted by date DESC")
	sort.Sort(sortByTime(files))

	i = 1
	for _, file := range files {
		size := float64(file.Metadata.Size) / 1024 / 1024

		fmt.Printf("File (created at: %s): %s = %f MB \n", file.Metadata.ClientModified.Format("2006-01-02"), file.Metadata.PathDisplay, size)

		if i >= SHOWRESULTS {
			break
		}

		i++
	}
}
