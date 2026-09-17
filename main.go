package main

import (
	"log"

	"github.com/AlexanderTal/catalog-service/internal/app/config"
	rhealth "github.com/AlexanderTal/catalog-service/internal/app/handler/http/health"
	rprocessor "github.com/AlexanderTal/catalog-service/internal/app/processor/http"
)

func main() {
	config.Load()

	cfg := config.Root

	hHealth := rhealth.NewHandler()

	httpServer := rprocessor.NewHTTP(hHealth, cfg.Processor.WebServer)
	if err := httpServer.Serve(); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
