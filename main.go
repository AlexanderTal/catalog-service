package main

import (
	"context"
	"log"

	"github.com/AlexanderTal/catalog-service/internal/app/config"
	hcategory "github.com/AlexanderTal/catalog-service/internal/app/handler/http/category"
	rhealth "github.com/AlexanderTal/catalog-service/internal/app/handler/http/health"
	hproduct "github.com/AlexanderTal/catalog-service/internal/app/handler/http/product"
	rprocessor "github.com/AlexanderTal/catalog-service/internal/app/processor/http"
	pcategory "github.com/AlexanderTal/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/AlexanderTal/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/AlexanderTal/catalog-service/internal/app/repository/product"
	scategory "github.com/AlexanderTal/catalog-service/internal/app/service/category"
	sproduct "github.com/AlexanderTal/catalog-service/internal/app/service/product"
)

func main() {
	ctx := context.Background()
	config.Load()
	cfg := config.Root

	pgClient, err := rcpostgres.NewClient(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	oldVer, newVer, err := pgClient.Migrate(ctx)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	if oldVer != newVer {
		log.Printf("Database migrated old_version=%d new_version=%d", oldVer, newVer)
	} else {
		log.Printf("Database is up to date version=%d", newVer)
	}

	// Репозитории
	categoryRepo := pcategory.NewRepoFromPostgres(pgClient)
	productRepo := pproduct.NewRepoFromPostgres(pgClient)

	// Сервисы
	categorySvc := scategory.NewService(categoryRepo, productRepo)
	productSvc := sproduct.NewService(productRepo, categoryRepo)

	// Хендлеры
	hHealth := rhealth.NewHandler()
	hCategory := hcategory.NewHandler(categorySvc)
	hProduct := hproduct.NewHandler(productSvc)

	// HTTP-сервер
	httpServer := rprocessor.NewHTTP(hHealth, hCategory, hProduct, cfg.Processor.WebServer)
	if err := httpServer.Serve(); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
