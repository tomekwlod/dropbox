package files

import (
	dropbox "github.com/tj/go-dropbox"
	config "github.com/tomekwlod/dropboxCleaner/config"
)

const perpage = 1000

func Search(term string) []*dropbox.SearchMatch {
	c := config.DropboxClient()

	var files []*dropbox.SearchMatch
	var i, start uint64 = 0, 0

	for {
		start = i * perpage

		out, err := c.Files.Search(&dropbox.SearchInput{
			Path:       "/",
			Query:      term,
			MaxResults: perpage,
			Start:      start,
		})

		if err != nil {
			panic(err)
		}

		for _, file := range out.Matches {
			files = append(files, file)
		}

		if !out.More {
			break
		}

		i++
	}

	return files
}

func Delete(term string) bool {
	c := config.DropboxClient()

	_, err := c.Files.Delete(&dropbox.DeleteInput{
		Path: term,
	})

	if err != nil {
		return false
	}

	return true
}
