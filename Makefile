.PHONY: help envs run migrate-up migrate-down migrate-new
.DEFAULT_GOAL := help
.SILENT: envs

DATABASE_CONN_URL ?= "postgres://postgres:password@localhost:5432/finanger?sslmode=disable"

help:
	echo "Usage: make [target]"
	echo ""
	echo "Targets:"
	echo "  envs          Set environment variables"
	echo "  run           Run the application"
	echo "  migrate-up    Run database migrations"
	echo "  migrate-down  Rollback database migrations"
	echo "  migrate-new   Create new database migration --> make migrate-new name=create_users_table"
	echo "  help          Show this help message"

run:
	go run cmd/main.go

envs-export:
	. ./.envs/local.env

migrate-up:
	migrate -path ./internal/app/database/migrations -database ${DATABASE_CONN_URL} up

migrate-down:
	migrate -path ./internal/app/database/migrations -database ${DATABASE_CONN_URL} down

migrate-new:
	migrate create -ext sql -dir ./internal/app/database/migrations ${name}