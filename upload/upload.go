package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	dropbox "github.com/tomekwlod/dropbox"
)

func main() {
	if len(os.Args[1:]) <= 1 {
		log.Println("You need to pass two parameters at least, eg: `command.go /target/on/dropbox file1.ext file2.ext`")

		return
	}

	for _, fn := range os.Args[2:] {
		in := &dropbox.UploadInput{
			Path:       os.Args[1] + "/" + filepath.Base(fn),
			Mode:       "add",
			Mute:       true,
			AutoRename: true,
		}

		result := dropbox.Upload(in)

		fmt.Printf("%+v", result)
	}

	log.Println("All done")
}
