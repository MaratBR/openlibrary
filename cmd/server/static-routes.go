package main

import (
	"net/http"

	uiserver "github.com/MaratBR/openlibrary/cmd/server/ui-server"
)

var devFrontEndServerProxy = uiserver.NewDevServerProxy(uiserver.DevServerOptions{
	GetServerPushedData: getServerData,
})

func devProxyIndex(w http.ResponseWriter, r *http.Request) {
	devFrontEndServerProxy.ServeHTTP(w, r)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
}
