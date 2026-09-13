# =============================================================================
# Переменные
# =============================================================================
OUTPUT := ./bin/app
GO_LINT_VERSION := 2.7.2
GO_FILE := ./main.go

# =============================================================================
# Справка
# =============================================================================
.PHONY: help
help: ## Показать справку
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# =============================================================================
# Разработка
# =============================================================================
.PHONY: run
run: ## Запустить приложение
	go run ${GO_FILE}

.PHONY: build
build: ## Сборка приложения
	go build -o ${OUTPUT} ${GO_FILE}

.PHONY: test
test: ## Запуск тестов
	go test -count=1 -v ./...

# =============================================================================
# Качество кода
# =============================================================================
.PHONY: lint
lint: ## Запуск линтера
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${GO_LINT_VERSION} run

.PHONY: lint-fix
lint-fix: ## Запуск линтера с автофиксом
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${GO_LINT_VERSION} run --fix

# =============================================================================
# Окружение (Docker)
# =============================================================================
.PHONY: up
up: ## Поднять docker окружение (PostgreSQL)
	docker compose up -d

.PHONY: down
down: ## Остановить docker окружение
	docker compose down --remove-orphans

.PHONY: logs
logs: ## Показать логи docker контейнеров
	docker compose logs -f

# =============================================================================
# Зависимости
# =============================================================================
.PHONY: deps
deps: ## Загрузить зависимости
	go mod tidy
	go mod download

.PHONY: mod-check
mod-check: ## Проверка актуальности go.mod/go.sum
	go mod tidy
	@FILES="go.mod"; [ -f go.sum ] && FILES="$$FILES go.sum"; git diff --exit-code -- $$FILES || (echo "go.mod/go.sum не синхронизированы. Запустите 'go mod tidy'" && exit 1)

# =============================================================================
# CI
# =============================================================================
.PHONY: ci
ci: ## Запустить все CI проверки
	@echo "=== Mod Check ==="
	go mod tidy
	@FILES="go.mod"; [ -f go.sum ] && FILES="$$FILES go.sum"; git diff --exit-code -- $$FILES || (echo "go.mod/go.sum не синхронизированы" && exit 1)
	@echo ""
	@echo "=== Build ==="
	@mkdir -p ./bin
	go build -o ./bin/ -v ./...
	@echo ""
	@echo "=== Test ==="
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo ""
	@echo "=== Lint ==="
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${GO_LINT_VERSION} run --timeout=10m
	@echo ""
	@echo "CI passed!"
