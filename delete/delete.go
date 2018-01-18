package main

import (
	"log"
	"os"

	"github.com/tomekwlod/dropboxCleaner/files"
)

func main() {
	fileNames := os.Args[1:]

	if len(fileNames) == 0 {
		log.Println("No files provided. Nothing to do")
		return
	}

	for _, fn := range fileNames {
		ok := files.Delete(fn)

		if ok == true {
			log.Printf("%v - Delete acknowledged", fn)
		} else {
			log.Printf("%v - No path found", fn)
		}
	}

}
