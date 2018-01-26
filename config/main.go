package config

import (
	dropbox "github.com/tj/go-dropbox"
	"github.com/tomekwlod/utils"
)

func DropboxToken() (token string, err error) {
	token, err = utils.EnvVariable("DROPBOX_ACCESS_TOKEN")
	return
}

func DropboxClient() *dropbox.Client {
	token, err := DropboxToken()

	if err != nil {
		panic("No token provided. Run `export DROPBOX_ACCESS_TOKEN=arUfR.......Cc`")
	}

	return dropbox.New(dropbox.NewConfig(token))
}
