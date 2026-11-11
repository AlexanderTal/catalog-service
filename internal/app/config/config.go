package config

import (
	"io"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/AlexanderTal/catalog-service/internal/app/config/section"
)

type Config struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}

var Root Config

func Load(args LoadArgs) {
	// Этап 1: Предварительный логгер с Debug уровнем

	// Переименование стандартных полей zerolog:
	// "time" -> "timestamp", "message" -> "msg"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC3339

	if args.EnableSimpleLog {
		args.Output = zerolog.ConsoleWriter{Out: args.Output}
	}

	// log.Logger — глобальная переменная пакета zerolog/log.
	log.Logger = createLogger(zerolog.DebugLevel, args.Output)

	log.Debug().Msg("Logger initialized with Debug level")

	// Загрузка конфигурации
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("No .env file found")
	}

	if err := envconfig.Process("APP", &Root); err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Этап 2: Переинициализация с уровнем из конфига
	level, err := zerolog.ParseLevel(Root.Monitor.LogLevel)
	if err != nil {
		// Если в конфиге опечатка, ставим Info по умолчанию и ругаемся в лог
		level = zerolog.InfoLevel
		log.Warn().Err(err).Str("log_level", Root.Monitor.LogLevel).Msg("Invalid log level, defaulting to info")
	}

	log.Logger = createLogger(level, args.Output)
	log.Info().Str("log_level", level.String()).Msg("Logger re-initialized with config level")
}

type LoadArgs struct {
	Output          io.Writer `json:"-"`
	EnableSimpleLog bool
}

func createLogger(level zerolog.Level, output io.Writer) zerolog.Logger {
	return zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()
}
