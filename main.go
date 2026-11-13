package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/AlexanderTal/catalog-service/cmd"
)

func main() {
	app := &cli.App{
		Name:    "catalog-service",
		Version: "1.0.0",
		Usage:   "Catalog management service",
		Commands: []*cli.Command{
			cmd.Migrate(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-json",
				Usage: "Human-readable log format",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
