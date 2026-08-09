//go:generate go run ../guts-maniac -output ../../web/frontend/src/api/generated.ts

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/MaratBR/openlibrary/cmd/server/util/minifycontent"
	"github.com/MaratBR/openlibrary/cmd/server/util/populate"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

func main() {
	log := newRootLogger()
	zap.ReplaceGlobals(log)
	defer func() { _ = log.Sync() }()

	command := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "util",
				Usage: "database maintenance and development utilities",
				Commands: []*cli.Command{
					{
						Name:  "minify-content",
						Usage: "sanitizes and minifies book summaries and chapter content",
						Action: func(ctx context.Context, c *cli.Command) error {
							config := loadConfigOrPanic()
							db := connectToDatabase(config, zap.S())
							if closer, ok := db.(interface{ Close() }); ok {
								defer closer.Close()
							}
							return minifycontent.Run(ctx, db, zap.S())
						},
					},
					{
						Name:  "populate",
						Usage: "populates the database with random data",
						Action: func(ctx context.Context, c *cli.Command) error {
							config := loadConfigOrPanic()
							db := connectToDatabase(config, zap.S())
							if closer, ok := db.(interface{ Close() }); ok {
								defer closer.Close()
							}
							return populate.Run(config, db, zap.S())
						},
					},
				},
			},
			{
				Name:  "server",
				Usage: "runs openlibrary server",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "dev",
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
					var cliParam cliParams
					cliParam.BypassTLSCheck = c.Bool("bypass-tls-check")
					cliParam.Dev = c.Bool("dev")
					cliParam.StaticDir = c.String("static-dir")
					mainServer(cliParam)
					return nil
				},
			},
		},
	}

	if err := command.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
