SHELL := /bin/bash
.DEFAULT_GOAL := help

BINARY  := bin/wsm
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X github.com/thasso/wsm/internal/cli.Version=$(VERSION)"

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the wsm binary
	@mkdir -p bin
	go build $(LDFLAGS) -o $(BINARY) ./cmd/wsm

.PHONY: install
install: ## Install wsm to $GOPATH/bin
	go install $(LDFLAGS) ./cmd/wsm

.PHONY: test
test: ## Run all tests
	go test ./...

.PHONY: fmt
fmt: ## Format Go source files
	gofmt -w .

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin/
