package main

import "github.com/go-chi/chi/v5"

type server struct {
	dev             bool
	disableTLSCheck bool
	r               chi.Router
}

func newServer(params *cliParams) *server {
	srv := &server{
		dev:             params.Dev,
		disableTLSCheck: params.BypassTLSCheck,
		r:               chi.NewRouter(),
	}

	return srv
}
