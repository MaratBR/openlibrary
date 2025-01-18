package frontend

import "net/http"

type AssetsConfig struct {
	Dev bool
}

func Assets(config AssetsConfig) http.Handler {
	if config.Dev {
		fileServer := http.FileServer(http.Dir("./web/frontend/dist"))
		return fileServer
	} else {
		panic("not implemented")
	}
}
