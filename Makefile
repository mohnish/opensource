.DEFAULT_GOAL := help

BINARY := opensource
VERSION ?= dev

.PHONY: help build install test vet fmt fmt-check version clean snapshot

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-12s %s\n", $$1, $$2}'

build: ## Build the binary into ./bin
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/$(BINARY) .

install: ## Install the binary into GOBIN
	go install -ldflags "-s -w -X main.version=$(VERSION)" .

test: ## Run the test suite
	go test -race ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Format the code
	gofmt -w .

fmt-check: ## Fail if any file needs formatting
	@test -z "$$(gofmt -l .)" || (echo "Run 'make fmt'"; gofmt -l .; exit 1)

version: ## Print the CLI version
	go run . --version

clean: ## Remove build artifacts
	rm -rf bin dist

snapshot: ## Build a local release snapshot with GoReleaser (no publish)
	go run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean
