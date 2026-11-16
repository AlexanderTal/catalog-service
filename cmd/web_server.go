package cmd

import (
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/AlexanderTal/catalog-service/internal/app/builder"
)

func WebServer() *cli.Command {
	return &cli.Command{
		Name:    "web-server",
		Aliases: []string{"ws"},
		Usage:   "Start HTTP server with all routes",
		Description: strings.TrimSpace(`
Initializes all dependencies (config, DB, repositories, services, handlers)
and starts the HTTP server. Graceful shutdown on SIGINT/SIGTERM.
`),
		Action:          cmdWebServer,
		HideHelpCommand: true,
	}
}

func cmdWebServer(cCtx *cli.Context) error {
	b := builder.NewBuilder(cCtx)
	b.BuildConfig()
	b.BuildRepoConnPostgres()
	b.BuildRepoConnMigrator()
	b.BuildRepoCategory()
	b.BuildRepoProduct()
	b.BuildServiceCategory()
	b.BuildServiceProduct()
	b.BuildHandlerHttpCategory()
	b.BuildHandlerHttpProduct()
	b.BuildProcHttp()
	b.Run()
	return nil
}
