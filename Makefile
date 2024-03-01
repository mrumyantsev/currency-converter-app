.SILENT:

include ./.env

.PHONY: fast-run
fast-run:
	go run ./cmd/${SERVER_APP_NAME}/main.go

.PHONY: build
build:
	go build -o ./build/${SERVER_APP_NAME} ./cmd/${SERVER_APP_NAME}/main.go

.PHONY: run
run:
	./build/${SERVER_APP_NAME}

.PHONY: save
save:
	./build/${SERVER_APP_NAME} -s

.PHONY: migrate
migrate: run-dbc migrate-up stop-dbc

.PHONY: run-dbc
run-dbc:
	docker run \
	--rm \
	--name ${DB_MIGRATION_CONTAINER_NAME} \
	-d \
	-e "PGDATA=/data" \
	-e "POSTGRES_USER=${DB_USERNAME}" \
	-e "POSTGRES_PASSWORD=${DB_PASSWORD}" \
	-e "POSTGRES_DB=${DB_DATABASE}" \
	-v ${DB_LOCAL_DIR}:/data \
	-p ${DB_MIGRATION_PORT}:5432 \
	postgres:${POSTGRES_VER}-alpine${ALPINE_VER}

.PHONY: migrate-up
migrate-up:
	./wait-for-postgres.sh \
	"${DB_PASSWORD}" \
	${DB_HOSTNAME} \
	${DB_MIGRATION_PORT} \
	${DB_USERNAME} \
	migrate \
	-path ./schema \
	-database "${DB_DRIVER}://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOSTNAME}:${DB_MIGRATION_PORT}/${DB_DATABASE}?sslmode=${DB_SSLMODE}" \
	up

.PHONY: stop-dbc
stop-dbc:
	docker stop ${DB_MIGRATION_CONTAINER_NAME}
