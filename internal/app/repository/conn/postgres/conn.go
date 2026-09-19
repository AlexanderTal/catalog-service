package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

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
	// Шаг 1: DSN
	var u url.URL
	u.Scheme = "postgres"
	u.Host = cfg.Address
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	u.Path = cfg.Name

	args := make(url.Values)
	args.Set("sslmode", "disable")
	u.RawQuery = args.Encode()

	dsn := u.String()

	log.Printf("Postgres timeouts: read=%s write=%s", cfg.ReadTimeout, cfg.WriteTimeout)

	// Шаг 2: sql.DB
	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithReadTimeout(cfg.ReadTimeout),
		pgdriver.WithWriteTimeout(cfg.WriteTimeout),
	))
	sqlDB.SetMaxOpenConns(10)

	// Шаг 3: bun.DB
	bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	// Шаг 4: Ping с таймаутом 2 секунды
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := bunDB.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	// Шаг 5: Client
	return &Client{
		_bunDB:   bunDB,
		rawBunDB: bunDB,
		cfg:      cfg,
	}, nil
}
func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	// Шаг 1: Загружаем миграции из embed
	migrations := migrate.NewMigrations()

	if err = migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	// Шаг 2: Создаём мигратор с опциями
	opts := []migrate.MigratorOption{
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable + "_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	}

	m := migrate.NewMigrator(c.rawBunDB, migrations, opts...)

	// Шаг 3: Инициализируем таблицу миграций
	if err = m.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("failed to init migration table: %w", err)
	}

	// Шаг 4: Получаем текущую версию (ДО применения)
	applied, err := m.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, mg := range applied {
		v, _ := strconv.ParseInt(mg.Name, 10, 64)
		if v > oldVer {
			oldVer = v
		}
	}

	// Шаг 5: Применяем новые миграции
	newMigrations, err := m.Migrate(ctx)
	if err != nil {
		return oldVer, oldVer, fmt.Errorf("failed to migrate: %w", err)
	}

	// Шаг 6: Вычисляем новую версию
	newVer = oldVer
	if newMigrations != nil {
		for _, mg := range newMigrations.Migrations {
			v, _ := strconv.ParseInt(mg.Name, 10, 64)
			if v > newVer {
				newVer = v
			}
		}
	}

	return oldVer, newVer, nil
}
