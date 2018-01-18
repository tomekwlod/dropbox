package config

import (
	env "github.com/segmentio/go-env"
	dropbox "github.com/tj/go-dropbox"
)

func DropboxClient() *dropbox.Client {
	token, err := env.Get("DROPBOX_ACCESS_TOKEN")

	if err != nil {
		panic("No token provided. Run `export DROPBOX_ACCESS_TOKEN=arUfR.......Cc`")
	}

	return dropbox.New(dropbox.NewConfig(token))
}
