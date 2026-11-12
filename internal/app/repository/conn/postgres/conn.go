package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/AlexanderTal/catalog-service/internal/app/config/section"
	"github.com/AlexanderTal/catalog-service/migration"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB

		cfg section.RepositoryPostgres
	}

	_bunDB = bun.IDB
)

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func NewClient(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {

	var u url.URL
	u.Scheme = "postgres"
	u.Host = cfg.Address
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	u.Path = cfg.Name

	args := make(url.Values)
	args.Set("sslmode", "disable")
	u.RawQuery = args.Encode()

	dsn := u.String()

	log.Info().Str("read_timeout", cfg.ReadTimeout.String()).Str("write_timeout", cfg.WriteTimeout.String()).Msg("Initializing PostgreSQL connection")

	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithReadTimeout(cfg.ReadTimeout),
		pgdriver.WithWriteTimeout(cfg.WriteTimeout),
	))
	sqlDB.SetMaxOpenConns(10)

	bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := bunDB.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &Client{
		_bunDB:   newTxInjector(bunDB),
		rawBunDB: bunDB,
		cfg:      cfg,
	}, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {

	migrations := migrate.NewMigrations()

	if err = migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	opts := []migrate.MigratorOption{
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable + "_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	}

	m := migrate.NewMigrator(c.rawBunDB, migrations, opts...)

	if err = m.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("failed to init migration table: %w", err)
	}

	applied, err := m.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, mg := range applied {
		v, err := strconv.ParseInt(mg.Name, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse migration version %q: %w", mg.Name, err)
		}
		if v > oldVer {
			oldVer = v
		}
	}

	newMigrations, err := m.Migrate(ctx)
	if err != nil {
		return oldVer, oldVer, fmt.Errorf("failed to migrate: %w", err)
	}

	newVer = oldVer
	if newMigrations != nil {
		for _, mg := range newMigrations.Migrations {
			v, err := strconv.ParseInt(mg.Name, 10, 64)
			if err != nil {
				return oldVer, oldVer, fmt.Errorf("failed to parse migration version %q: %w", mg.Name, err)
			}
			if v > newVer {
				newVer = v
			}
		}
	}

	return oldVer, newVer, nil
}

func (c *Client) InsideTx(ctx context.Context, fn func(ctx context.Context) error) error {

	if tx := getTxFromContext(ctx); tx.Tx != nil {
		return fn(ctx)
	}

	tx, err := c.rawBunDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	done := false
	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	txCtx := setTxToContext(ctx, tx)
	if err := fn(txCtx); err != nil {
		return err // defer сделает Rollback
	}

	done = true
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}
