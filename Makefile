.DEFAULT_GOAL := build

OUT_DIR ?= build
EXE ?= $(OUT_DIR)/openlibrary-server

LOCAL_DB_HOST ?= localhost
LOCAL_DB_USER ?= postgres
LOCAL_DB_PASSWORD ?= postgres
LOCAL_DB_PORT ?= 5432
LOCAL_DB_NAME ?= openlibrary
LOCAL_DB ?= postgres://$(LOCAL_DB_USER):$(LOCAL_DB_PASSWORD)@$(LOCAL_DB_HOST):$(LOCAL_DB_PORT)/$(LOCAL_DB_NAME)?sslmode=disable
PGX_MIGRATIONS ?= file://internal/store/migrations
MIGRATE_ARGS = -source=$(PGX_MIGRATIONS) -database=$(LOCAL_DB)

.PHONY: help verify_go verify_gow verify_templ verify_migrate verify_psql \
	verify_sqlc verify_pnpm verify_docker build build_server main_watch templ templ_watch \
	migration migrate_db db_reset db_migrate_down_1 db_sqlc ui_watch ui_build \
	codegen check_generated test check db_populate db_backup

help:
	@echo "Common targets:"
	@echo "  build             Build the Go server"
	@echo "  codegen           Generate SQLC and templ output"
	@echo "  check             Generate, test, build the UI, and check whitespace"
	@echo "  main_watch        Run the development server with gow"
	@echo "  ui_watch          Build frontend assets in watch mode"
	@echo "  migration N=name  Create a database migration"
	@echo "  migrate_db        Apply all pending database migrations"
	@echo "  db_reset CONFIRM=yes  Recreate and migrate the local database"

verify_go:
	@./scripts/verify_cli.sh go

verify_gow:
	@./scripts/verify_cli.sh gow

verify_templ:
	@./scripts/verify_cli.sh templ

verify_migrate:
	@./scripts/verify_cli.sh migrate

verify_psql:
	@./scripts/verify_cli.sh psql

verify_sqlc:
	@./scripts/verify_cli.sh sqlc

verify_pnpm:
	@./scripts/verify_cli.sh pnpm

verify_docker:
	@./scripts/verify_cli.sh docker

build_server: verify_go
	@mkdir -p "$(OUT_DIR)"
	go build -o "$(EXE)" ./cmd/server

build: build_server

main_watch: verify_gow
	gow run ./cmd/server server --dev --bypass-tls-check --static-dir ./cmd/server/ui/dist

templ: verify_templ
	templ generate

templ_watch: verify_templ
	templ generate --watch

migration: verify_migrate
	@test -n "$(N)" || { echo "Usage: make migration N=<name>"; exit 1; }
	migrate create -ext sql -dir internal/store/migrations -seq "$(N)"

migrate_db: verify_migrate
	migrate $(MIGRATE_ARGS) up

db_reset: verify_psql verify_migrate
	@test "$(CONFIRM)" = "yes" || { echo "This drops $(LOCAL_DB_NAME). Run: make db_reset CONFIRM=yes"; exit 1; }
	PGPASSWORD="$(LOCAL_DB_PASSWORD)" psql -p "$(LOCAL_DB_PORT)" -h "$(LOCAL_DB_HOST)" -U "$(LOCAL_DB_USER)" -c "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '$(LOCAL_DB_NAME)' AND pid <> pg_backend_pid();"
	PGPASSWORD="$(LOCAL_DB_PASSWORD)" psql -p "$(LOCAL_DB_PORT)" -h "$(LOCAL_DB_HOST)" -U "$(LOCAL_DB_USER)" -c "DROP DATABASE IF EXISTS \"$(LOCAL_DB_NAME)\""
	PGPASSWORD="$(LOCAL_DB_PASSWORD)" psql -p "$(LOCAL_DB_PORT)" -h "$(LOCAL_DB_HOST)" -U "$(LOCAL_DB_USER)" -c "CREATE DATABASE \"$(LOCAL_DB_NAME)\""
	migrate $(MIGRATE_ARGS) up

db_migrate_down_1: verify_migrate
	migrate $(MIGRATE_ARGS) down 1

db_sqlc: verify_sqlc
	sqlc -f internal/store/sqlc.yaml generate

ui_watch: verify_pnpm
	pnpm run dev

ui_build: verify_pnpm
	pnpm run build

codegen: db_sqlc templ

check_generated: codegen
	@git diff --exit-code -- 'internal/store/*.sql.go' internal/store/models.go internal/store/db.go internal/store/copyfrom.go || { echo "Generated SQLC files are out of date"; exit 1; }

test: verify_go
	go test ./...

check: check_generated test ui_build
	git diff --check

db_populate: verify_go
	go run ./cmd/server populate

db_backup: verify_docker
	@mkdir -p "$(OUT_DIR)"
	docker compose exec -T postgres pg_dump -U "$(LOCAL_DB_USER)" -F c -b -v "$(LOCAL_DB_NAME)" > "$(OUT_DIR)/openlibrary.backup"
