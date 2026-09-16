package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/AlexanderTal/catalog-service/internal/app/config/section"
)

type Config struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}

var Root Config

func Load() {
	_ = godotenv.Load()

	if err := envconfig.Process("APP", &Root); err != nil {
		log.Fatal(err)
	}
}
