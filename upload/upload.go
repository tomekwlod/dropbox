package main

import (
	"fmt"
	"log"
	"os"

	dropbox "github.com/tomekwlod/dropbox"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("You need to pass exactly two parameters, eg: `command file.ext /target/on/dropbox/file.ext`")

		return
	}

	in := &dropbox.UploadInput{
		Path:       os.Args[2],
		LocalFile:  os.Args[1],
		Mode:       "add",
		Mute:       true,
		AutoRename: true,
	}

	result := dropbox.Upload(in)

	fmt.Printf("%+v", result)

	log.Println("All done")
}
