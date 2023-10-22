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

test:
	go test -race ./internal/...

test-local:
	CALENDAR_DATABASE_DSN=$(DB_TEST_DSN) go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run --timeout=2m ./...

up-related-services:
	docker-compose -f ./local/docker-compose.yaml up -d

cleanup:
	docker-compose -f ./local/docker-compose.yaml down

migrate-local:
	GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_LOCAL_DSN) goose -dir $(DB_MIGRATIONS_DIR) $(direction)  # direction = up | up-by-one | down

install-pb-deps:
	(which protoc-gen-go > /dev/null) || go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	(which protoc-gen-go-grpc > /dev/null) || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	(which protoc-gen-validate > /dev/null) || go install github.com/envoyproxy/protoc-gen-validate@v1.0.2
	(which buf > /dev/null) || go install github.com/bufbuild/buf/cmd/buf@1.27.1

generate-pb: install-pb-deps
	cd ./api; buf generate

.PHONY: build run build-img run-img version test test-local lint up-related-services cleanup migrate-local generate-pb
