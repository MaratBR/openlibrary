package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MaratBR/openlibrary/cmd/go2tsdef/generator"
)

func main() {
	configPath := flag.String("config", "go2tsdef.yaml", "path to the YAML configuration file")
	flag.Parse()

	if err := generator.Run(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "go2tsdef: %v\n", err)
		os.Exit(1)
	}
}
