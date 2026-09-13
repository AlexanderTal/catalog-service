# MoM Boilerplate

Стартовый шаблон для Go-микросервисов в рамках Battle Project.

Содержит только инфраструктурную обвязку — вся Go-логика реализуется в задачах курса.

## Требования

| Зависимость                  | Версия |
|------------------------------|--------|
| [Go](https://go.dev)         | 1.25+  |
| [Docker](https://docker.com) | 20+    |
| make                         | -      |

## Быстрый старт

### 1. Клонирование шаблона

```bash
git clone --depth=1 git@github.com:MoM-Repo/mom-boilerplate.git YOUR_PROJECT_NAME
cd YOUR_PROJECT_NAME
```

### 2. Настройка проекта

```bash
chmod +x setup.sh
./setup.sh your-project-name your-github-username
```

Скрипт автоматически:
- Обновит все импорты в Go файлах
- Изменит `go.mod` на новый модуль
- Выполнит `go mod tidy`

### 3. Инициализация репозитория

```bash
rm -rf .git
git init
git add .
git commit -m "Initial commit"
git remote add origin git@github.com:your-username/your-project-name.git
git push -u origin main
```

### 4. Запуск

```bash
# Поднять PostgreSQL
make up

# Запустить приложение
make run
```

## Доступные команды

```bash
make help  # Показать все команды
```

| Команда          | Описание                          |
|------------------|-----------------------------------|
| `make run`       | Запустить приложение              |
| `make build`     | Собрать бинарник                  |
| `make test`      | Запустить тесты                   |
| `make lint`      | Запустить линтер                  |
| `make lint-fix`  | Линтер + автофикс                 |
| `make up`        | Поднять PostgreSQL                |
| `make down`      | Остановить контейнеры             |
| `make logs`      | Показать логи контейнеров         |
| `make deps`      | Загрузить зависимости             |
| `make mod-check` | Проверка актуальности go.mod      |
| `make ci`        | Запустить все CI проверки локально |

## Структура проекта

```
.
├── .github/workflows/      # GitHub Actions CI
├── cmd/                    # CLI команды (urfave/cli)
├── internal/
│   └── app/
│       ├── builder/        # Dependency injection
│       ├── config/         # Конфигурация
│       ├── entity/         # Модели/DTO
│       ├── handler/        # HTTP/gRPC обработчики
│       ├── processor/      # Процессоры (HTTP server, etc.)
│       ├── repository/     # Слой данных
│       ├── service/        # Бизнес-логика
│       └── util/           # Утилиты
├── docker-compose.yml      # PostgreSQL для разработки
├── .golangci.yml           # Конфигурация линтера (v2)
├── Makefile
└── setup.sh                # Скрипт переименования проекта
```
