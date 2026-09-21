package main

import (
	"context"
	"log"

	"github.com/AlexanderTal/catalog-service/internal/app/config"
	rhealth "github.com/AlexanderTal/catalog-service/internal/app/handler/http/health"
	rprocessor "github.com/AlexanderTal/catalog-service/internal/app/processor/http"
	rcpostgres "github.com/AlexanderTal/catalog-service/internal/app/repository/conn/postgres"
)

func main() {
	ctx := context.Background()
	config.Load()
	cfg := config.Root

	// Подключение к PostgreSQL
	pgClient, err := rcpostgres.NewClient(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Применение миграций
	oldVer, newVer, err := pgClient.Migrate(ctx)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if oldVer != newVer {
		log.Printf("Database migrated old_version=%d new_version=%d", oldVer, newVer)
	} else {
		log.Printf("Database is up to date version=%d", newVer)
	}

	// HTTP-сервер (код из TASK-1002)
	hHealth := rhealth.NewHandler()

	httpServer := rprocessor.NewHTTP(hHealth, cfg.Processor.WebServer)
	if err := httpServer.Serve(); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
