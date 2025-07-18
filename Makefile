OUT_DIR := build
EXE := $(OUT_DIR)/openlibrary-server

LOCAL_DB_HOST := localhost
LOCAL_DB_USER := postgres
LOCAL_DB_PASSWORD := postgres
LOCAL_DB_PORT := 5432
LOCAL_DB := postgres://$(LOCAL_DB_USER):$(LOCAL_DB_PASSWORD)@$(LOCAL_DB_HOST):$(LOCAL_DB_PORT)/openlibrary?sslmode=disable
PGX_MIGRATIONS := file://internal/store/migrations
MIGRATE_ARGS := -source=$(PGX_MIGRATIONS) -database=$(LOCAL_DB)

build-server:
	go build -o $(EXE) ./cmd/server

build: sqlc templ build-server

watch-server:
	gow run ./cmd/server server --dev --bypass-tls-check --static-dir ./cmd/server/ui/dist

watch-templ:
	templ generate --watch


migration:
	migrate create -ext sql -dir internal/store/migrations -seq $N

migrate-db:
	migrate -source=$(PGX_MIGRATIONS) -database=$(LOCAL_DB) up

reset-db:
	PGPASSWORD=$(LOCAL_DB_PASSWORD) psql -p $(LOCAL_DB_PORT) -h $(LOCAL_DB_HOST) -U $(LOCAL_DB_USER) -c "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = 'openlibrary' AND pid <> pg_backend_pid();"
	PGPASSWORD=$(LOCAL_DB_PASSWORD) psql -p $(LOCAL_DB_PORT) -h $(LOCAL_DB_HOST) -U $(LOCAL_DB_USER) -c "DROP DATABASE IF EXISTS openlibrary"
	PGPASSWORD=$(LOCAL_DB_PASSWORD) psql -p $(LOCAL_DB_PORT) -h $(LOCAL_DB_HOST) -U $(LOCAL_DB_USER) -c "CREATE DATABASE openlibrary"
	migrate $(MIGRATE_ARGS) up

migrate-db-down-1:
	migrate $(MIGRATE_ARGS) down 1 

templ:
	templ generate

sqlc:
	sqlc -f internal/store/sqlc.yaml generate


proto:
	mkdir -p ./cmd/server/front/src/proto

	# search protobuf
	protoc -I=./cmd/server/olproto --go_out=. ./cmd/server/olproto/search.proto
	protoc -I=./cmd/server/olproto \
		--plugin=./cmd/server/front/node_modules/ts-proto/protoc-gen-ts_proto \
		--ts_proto_out=./cmd/server/front/src/proto \
		--ts_proto_opt=forceLong=string \
		./cmd/server/olproto/search.proto

ao3-build-docker:
	sudo docker build -t openlibrary/ao3-scrapper -f ./cmd/ao3-scrapper/Dockerfile .

#
# FRONT-END
#

watch-front:
	cd ./web/frontend && pnpm run dev

build-front:
	cd ./web/frontend && pnpm run build