SHELL := /bin/sh

# Config
APP_NAME := ./bin/blog-api
IMAGE_NAME := blog-api
VERSION ?= latest
DATABASE_URL ?= postgres://postgres:postgres@localhost/blogdb?sslmode=disable
ENV_FILE ?= config/.env

# Helper to load env file into the shell session for each recipe
with_env = set -a; [ -f $(ENV_FILE) ] && . $(ENV_FILE); set +a;

# Migration tool: golang-migrate (https://github.com/golang-migrate/migrate)
MIGRATE := migrate
MIGR_DIR := migrations

.PHONY: help
help:
	@echo "Targets:"
	@echo "  migrate-create name=...   Create a new migration"
	@echo "  migrate-up                Apply all up migrations"
	@echo "  migrate-down              Rollback one migration"
	@echo "  build-http                Build host OS binary"
	@echo "  run-http                  Run the HTTP server"
	@echo "  test                      Run all tests with verbose output"
	@echo "  test-coverage             Run tests and generate coverage report"
	@echo "  swagger-ui                Launch Swagger UI to view docs/openapi.yaml"
	@echo "  docker-build [VERSION]    Build Docker image"
	@echo "  docker-save [VERSION]     Save Docker image to tar"

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=add_feature"; exit 1; fi
	$(MIGRATE) create -ext sql -dir $(MIGR_DIR) -seq $(name)

.PHONY: migrate-up
migrate-up:
	@$(with_env) $(MIGRATE) -path $(MIGR_DIR) -database "$$DATABASE_URL" up

.PHONY: migrate-down
migrate-down:
	@$(with_env) $(MIGRATE) -path $(MIGR_DIR) -database "$$DATABASE_URL" down 1

.PHONY: build-http
build-http:
	go build -o $(APP_NAME) ./cmd/http
	@echo executable file saved in $(APP_NAME)

.PHONY: run-http
run-http:
	./$(APP_NAME) --env-path="./config/.env"

.PHONY: docker-build
docker-build:
	@$(with_env) tag=$${VERSION:-$(VERSION)}; docker build -t $(IMAGE_NAME):$$tag .

.PHONY: docker-save
docker-save:
	@$(with_env) tag=$${VERSION:-$(VERSION)}; docker save -o $(IMAGE_NAME)_$$tag.tar $(IMAGE_NAME):$$tag

.PHONY: env-print
env-print:
	@$(with_env) env | sort | grep -E '^(DATABASE_URL|PORT|VERSION)='

.PHONY: test
test:
	go test ./... -v

.PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: swagger-ui
swagger-ui:
	@docker run -d --rm -p 8081:8080 -e SWAGGER_JSON=/docs/openapi.yaml -v $(PWD)/docs:/docs --name swagger-ui swaggerapi/swagger-ui



