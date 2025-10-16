OUT_DIR := build
EXE := $(OUT_DIR)/openlibrary-server

LOCAL_DB_HOST := localhost
LOCAL_DB_USER := postgres
LOCAL_DB_PASSWORD := postgres
LOCAL_DB_PORT := 5432
LOCAL_DB := postgres://$(LOCAL_DB_USER):$(LOCAL_DB_PASSWORD)@$(LOCAL_DB_HOST):$(LOCAL_DB_PORT)/openlibrary?sslmode=disable
PGX_MIGRATIONS := file://internal/store/migrations
MIGRATE_ARGS := -source=$(PGX_MIGRATIONS) -database=$(LOCAL_DB)

build_server:
	go build -o $(EXE) ./cmd/server

build: codegen build_server

main_watch:
	gow run ./cmd/server server --dev --bypass-tls-check --static-dir ./cmd/server/ui/dist

templ:
	templ generate


templ_watch:
	templ generate --watch


migration:
	migrate create -ext sql -dir internal/store/migrations -seq $N

migrate-db:
	migrate -source=$(PGX_MIGRATIONS) -database=$(LOCAL_DB) up

db_reset:
	PGPASSWORD=$(LOCAL_DB_PASSWORD) psql -p $(LOCAL_DB_PORT) -h $(LOCAL_DB_HOST) -U $(LOCAL_DB_USER) -c "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = 'openlibrary' AND pid <> pg_backend_pid();"
	PGPASSWORD=$(LOCAL_DB_PASSWORD) psql -p $(LOCAL_DB_PORT) -h $(LOCAL_DB_HOST) -U $(LOCAL_DB_USER) -c "DROP DATABASE IF EXISTS openlibrary"
	PGPASSWORD=$(LOCAL_DB_PASSWORD) psql -p $(LOCAL_DB_PORT) -h $(LOCAL_DB_HOST) -U $(LOCAL_DB_USER) -c "CREATE DATABASE openlibrary"
	migrate $(MIGRATE_ARGS) up

db_migrate_down_1:
	migrate $(MIGRATE_ARGS) down 1

db_sqlc:
	sqlc -f internal/store/sqlc.yaml generate


# proto:
# 	mkdir -p ./cmd/server/front/src/proto

# 	# search protobuf
# 	protoc -I=./cmd/server/olproto --go_out=. ./cmd/server/olproto/search.proto
# 	protoc -I=./cmd/server/olproto \
# 		--plugin=./cmd/server/front/node_modules/ts-proto/protoc-gen-ts_proto \
# 		--ts_proto_out=./cmd/server/front/src/proto \
# 		--ts_proto_opt=forceLong=string \
# 		./cmd/server/olproto/search.proto

# ao3-build-docker:
# 	sudo docker build -t openlibrary/ao3-scrapper -f ./cmd/ao3-scrapper/Dockerfile .

#
# FRONT-END
#

ui_watch:
	cd ./web/frontend && pnpm run dev

ui_build:
	cd ./web/frontend && pnpm run build

codegen: db_sqlc templ

db_populate:
	go run ./cmd/server populate

db_backup:
	sudo docker run --rm \
		-e PGPASSWORD=$(LOCAL_DB_PASSWORD) \
		-v ./$(OUT_DIR):/backup \
		--network host \
		postgres:latest \
		pg_dump -h $(LOCAL_DB_HOST) -p $(LOCAL_DB_PORT) -U $(LOCAL_DB_USER) \
		-F c -b -v -f "/backup/openlibrary.backup" openlibrary
