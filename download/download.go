package main

import (
	"log"
	"os"

	"github.com/tomekwlod/dropbox"
)

func main() {
	if len(os.Args) <= 1 {
		log.Println("You must pass exactly two parameters which tell what and where to download, eg: `command file/to/download.ext local/file.ext`")
	}

	filename := os.Args[1]
	target := os.Args[2]

	err := dropbox.Download(filename, target)
	if err != nil {
		panic(err)
	}
	return

}
