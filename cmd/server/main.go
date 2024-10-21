package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.NotFound(index)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
