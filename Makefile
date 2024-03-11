.SILENT:
.DEFAULT_GOAL := fast-run

include ./.env

export ALPINE_VER
export DB_DATABASE
export DB_DRIVER
export DB_HOSTNAME
export DB_LOCAL_DIR
export DB_MIGRATION_CONTAINER_NAME
export DB_MIGRATION_PORT
export DB_PASSWORD
export DB_PORT
export DB_SSLMODE
export DB_USERNAME
export ENABLE_DEBUG_LOGS
export GO_VER
export HTTP_SERVER_LISTEN_PORT
export NGINX_VER
export POSTGRES_VER
export SERVER_APP_NAME
export WEB_LOCAL_DIR
export WEB_PORT

.PHONY: migrate
migrate: run-dbc migrate-up stop-dbc

.PHONY: run-dbc
run-dbc:
	docker run \
	--rm \
	--name ${DB_MIGRATION_CONTAINER_NAME} \
	-d \
	-e PGDATA=/data \
	-e POSTGRES_USER=${DB_USERNAME} \
	-e POSTGRES_PASSWORD=${DB_PASSWORD} \
	-e POSTGRES_DB=${DB_DATABASE} \
	-v ${DB_LOCAL_DIR}:/data \
	-p ${DB_MIGRATION_PORT}:5432 \
	postgres:${POSTGRES_VER}-alpine${ALPINE_VER}

.PHONY: migrate-up
migrate-up:
	./wait-for-postgres.sh \
	${DB_PASSWORD} \
	${DB_HOSTNAME} \
	${DB_MIGRATION_PORT} \
	${DB_USERNAME} \
	migrate \
	-path ./schema \
	-database ${DB_DRIVER}://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOSTNAME}:${DB_MIGRATION_PORT}/${DB_DATABASE}?sslmode=${DB_SSLMODE} \
	up

.PHONY: stop-dbc
stop-dbc:
	docker stop ${DB_MIGRATION_CONTAINER_NAME}

.PHONY: build
build:
	go build -o ./build/${SERVER_APP_NAME} ./cmd/${SERVER_APP_NAME}/main.go

.PHONY: run
run:
	./build/${SERVER_APP_NAME}

.PHONY: fast-run
fast-run:
	go run ./cmd/${SERVER_APP_NAME}/main.go

.PHONY: save
save:
	./build/${SERVER_APP_NAME} -s
