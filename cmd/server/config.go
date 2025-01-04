package main

import (
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func loadConfigOrPanic() *koanf.Koanf {
	k := koanf.New(".")
	if err := k.Load(file.Provider("openlibrary.toml"), toml.Parser()); err != nil {
		panic(err)
	}
	return k
}
