package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	command := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "populate",
				Usage: "populates database with random data",
				Action: func(ctx context.Context, c *cli.Command) error {
					cfg := loadConfigOrPanic()
					mainPopulate(cfg)
					return nil
				},
			},
			{
				Name:  "server",
				Usage: "runs openlibrary server",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "dev-frontend-proxy",
						Usage: "enable dev frontend proxy",
					},
					&cli.BoolFlag{
						Name:  "bypass-tls-check",
						Usage: "disables TLS check when exchanging sensitive data, such as when user signs in or signs up and plain text password is being exchanged",
					},
					&cli.StringFlag{
						Name:  "static-dir",
						Usage: "directory with static files",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					cfg := loadConfigOrPanic()
					var cliParam cliParams
					cliParam.BypassTLSCheck = c.Bool("bypass-tls-check")
					cliParam.Dev = c.Bool("dev-frontend-proxy")
					cliParam.StaticDir = c.String("static-dir")
					mainServer(cliParam, cfg)
					return nil
				},
			},
		},
	}

	command.Run(context.Background(), os.Args)
}
