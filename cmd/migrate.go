package cmd

import (
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/AlexanderTal/catalog-service/internal/app/builder"
)

func Migrate() *cli.Command {
	return &cli.Command{
		Name:    "migrate",
		Aliases: []string{"m"},
		Usage:   "Apply pending database migrations",
		Description: strings.TrimSpace(`
Connects to PostgreSQL, checks current schema version,
and applies any pending migrations.
`),
		Action:          cmdMigrate,
		HideHelpCommand: true,
	}
}

func cmdMigrate(cCtx *cli.Context) error {
	b := builder.NewBuilder(cCtx)
	b.BuildConfig()
	b.BuildRepoConnPostgres()
	b.BuildRepoConnMigrator()
	b.Run()
	return nil
}
