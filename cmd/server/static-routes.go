package main

import (
	"net/http"

	uiserver "github.com/MaratBR/openlibrary/cmd/server/ui-server"
)

var devFrontEndServerProxy = uiserver.NewDevServerProxy()

func index(w http.ResponseWriter, r *http.Request) {
	devFrontEndServerProxy.ServeHTTP(w, r)
}
