package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/tomekwlod/dropbox"
)

const SHOWRESULTS = 20

type sortBySize []*dropbox.FilesResult

func (a sortBySize) Len() int           { return len(a) }
func (a sortBySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortBySize) Less(i, j int) bool { return a[i].Size > a[j].Size }

func main() {
	limit := SHOWRESULTS
	var err error

	if len(os.Args) > 1 {
		limit, err = strconv.Atoi(os.Args[1])

		if err != nil {
			log.Printf("Oprional parameter in number type needed, %s given", os.Args[1])
		}
	}

	log.Printf("Getting all the files")

	files := dropbox.AllFiles()

	if len(files) == 0 {
		log.Printf("No files found")

		return
	}

	log.Printf("%d files found\n", len(files))

	fmt.Println("\nSorted by size DESC")
	sort.Sort(sortBySize(files))

	i := 1
	for _, file := range files {
		size := float64(file.Size) / 1024 / 1024

		fmt.Printf("%d]: File (created at: %s): %s = %f MB \n", i, file.Modified.Format("2006-01-02"), file.Name, size)

		if i >= limit {
			break
		}

		i++
	}
}
