package main

import (
	"flag"

	"github.com/MaratBR/openlibrary/cmd/server"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var cliParam server.CLIParams

func loadConfigOrPanic() *koanf.Koanf {
	k := koanf.New(".")
	if err := k.Load(file.Provider("openlibrary.toml"), toml.Parser()); err != nil {
		panic(err)
	}
	return k
}

func main() {
	flag.BoolVar(&cliParam.Dev, "dev-frontend-proxy", false, "enable dev frontend proxy")
	flag.BoolVar(&cliParam.BypassTLSCheck, "bypass-tls-check", false, "disables TLS check when exchanging sensitive data")

	flag.Parse()

	cfg := loadConfigOrPanic()

	server.Main(cliParam, cfg)
}
