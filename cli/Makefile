.PHONY: build clean fmt lint test test-cover run

BIN_DIR    := bin
BIN_NAME   := bangumi
ifeq ($(OS),Windows_NT)
	BIN_NAME := bangumi.exe
endif

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -s -w \
	-X 'main.version=$(VERSION)' \
	-X 'main.commit=$(COMMIT)' \
	-X 'main.date=$(DATE)' \
	-X 'github.com/kasuganosora/bangumi.skill/cli/api.version=$(VERSION)'

build:
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BIN_NAME) .

clean:
	rm -rf $(BIN_DIR)

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

test:
	go test ./... -v -count=1

test-cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

run: build
	./$(BIN_DIR)/$(BIN_NAME)
