package builder

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/AlexanderTal/catalog-service/internal/app/config"
	"github.com/AlexanderTal/catalog-service/internal/app/processor"
	pprocessor "github.com/AlexanderTal/catalog-service/internal/app/processor/other"
	"github.com/AlexanderTal/catalog-service/internal/app/repository"
	pcategory "github.com/AlexanderTal/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/AlexanderTal/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/AlexanderTal/catalog-service/internal/app/repository/product"
)

type Builder struct {
	cCtx *cli.Context
	ctx  context.Context
	wg   sync.WaitGroup
	err  error
	cfg  config.Config

	connPostgres *rcpostgres.Client
	categoryRepo repository.Category
	productRepo  repository.Product
	processors   []processor.Processor
}

func NewBuilder(cCtx *cli.Context) *Builder {
	return &Builder{
		cCtx: cCtx,
		ctx:  context.Background(),
	}
}

func (b *Builder) BuildConfig() {
	b.exec(func(b *Builder) {
		b.buildConfig()
	})
}

func (b *Builder) Run() {
	if b.err != nil {
		log.Fatal().Err(b.err).Msg("Failed to initialize application")
	}

	log.Info().Msg("Application initialized")
	defer log.Info().Msg("Application completed")

	for _, proc := range b.processors {
		proc.StartAsync(b.ctx, &b.wg)
	}

	b.wg.Wait()
}

////////////////////////////////////////////////////////////////////////////////
///// REPOSITORY CONNECTIONS ///////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildRepoConnPostgres() {
	b.exec(func(b *Builder) {
		client, err := rcpostgres.NewClient(b.ctx, b.cfg.Repository.Postgres)
		if err != nil {
			b.err = err
			return
		}
		b.connPostgres = client
	})
}

func (b *Builder) BuildRepoConnMigrator() {
	if b.connPostgres == nil {
		return
	}

	b.exec(func(b *Builder) {
		migrator := pprocessor.NewMigrator(b.connPostgres)
		b.processors = append(b.processors, migrator)
	})
}

////////////////////////////////////////////////////////////////////////////////
///// REPOSITORIES /////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildRepoCategory() {
	b.exec(func(b *Builder) {
		b.categoryRepo = pcategory.NewRepoFromPostgres(b.connPostgres)
	}, b.connPostgres)
}

func (b *Builder) BuildRepoProduct() {
	b.exec(func(b *Builder) {
		b.productRepo = pproduct.NewRepoFromPostgres(b.connPostgres)
	}, b.connPostgres)
}

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE //////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) buildConfig() {
	args := config.LoadArgs{
		Output:          b.cCtx.App.Writer,
		EnableSimpleLog: b.cCtx.Bool("no-json"),
	}
	config.Load(args)
	b.cfg = config.Root
}

func (b *Builder) exec(cb func(b *Builder), requiredArgs ...any) {
	if b.err != nil {
		return
	}

	for i, requiredArg := range requiredArgs {
		rv := reflect.ValueOf(requiredArg)
		if !rv.IsValid() {
			b.err = fmt.Errorf("BUG: required argument #%d is nil (check dependencies)", i)
			return
		}
		if rv.Type().Kind() == reflect.Struct || !rv.IsZero() {
			continue
		}
		b.err = fmt.Errorf("BUG: required %s, but empty", rv.Type().String())
		return
	}

	cb(b)
}
