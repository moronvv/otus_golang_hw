BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DB_DRIVER=postgres
DB_MIGRATIONS_DIR=./migrations
DB_LOCAL_DSN="host=localhost user=user password=dbpass dbname=calendar sslmode=disable"
DB_TEST_DSN="host=localhost user=user password=dbpass dbname=test sslmode=disable"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) --config ./configs/calendar.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) --version

test-local:
	CALENDAR_DATABASE_DSN=$(DB_TEST_DSN) go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

up-related-services:
	docker-compose -f ./local/docker-compose.yaml up -d

cleanup:
	docker-compose -f ./local/docker-compose.yaml down

migrate-local:
	GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_LOCAL_DSN) goose -dir $(DB_MIGRATIONS_DIR) $(direction)  # direction = up | up-by-one | down

.PHONY: build run build-img run-img version test-local lint up-related-services cleanup migrate-local
