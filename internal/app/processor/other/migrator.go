package pprocessor

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/AlexanderTal/catalog-service/internal/app/processor"
	"github.com/AlexanderTal/catalog-service/internal/app/repository"
)

type procMigrate struct {
	migrator repository.Migrate
}

func NewMigrator(migrator repository.Migrate) processor.Processor {
	return &procMigrate{
		migrator: migrator,
	}
}

func (p *procMigrate) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	processor.Wrap(ctx, wg, p.job)
}

func (p *procMigrate) job(ctx context.Context) {
	oldVer, newVer, err := p.migrator.Migrate(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to apply migrations")
		return
	}

	if oldVer != newVer {
		log.Info().
			Int64("old_version", oldVer).
			Int64("new_version", newVer).
			Msg("Repository schema has been updated")
	} else {
		log.Info().Int64("version", newVer).Msg("Repository schema is up to date")
	}
}
