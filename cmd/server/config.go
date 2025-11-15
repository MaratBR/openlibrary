package main

import (
	"os"

	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func loadConfigOrPanic() *koanf.Koanf {
	k := koanf.New(".")
	if err := k.Load(file.Provider("openlibrary.toml"), toml.Parser()); err != nil {
		panic(err)
	}
	if _, err := os.Stat("openlibrary.private.toml"); err == nil {
		if err := k.Load(file.Provider("openlibrary.private.toml"), toml.Parser()); err != nil {
			panic(err)
		}
	}
	return k
}
