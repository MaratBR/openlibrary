package frontend

import (
	"io/fs"
	"net/http"
	"os"
)

type AssetsConfig struct {
	Dev bool
}

func AssetsFS(config AssetsConfig) fs.FS {
	if config.Dev {
		return os.DirFS("./dist")
	} else {
		panic("not implemented")
	}
}

func Assets(fs fs.FS) http.Handler {
	fileServer := http.FileServerFS(fs)
	return fileServer
}
