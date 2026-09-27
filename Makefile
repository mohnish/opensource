.DEFAULT_GOAL := help

BINARY := opensource
VERSION ?= dev
GORELEASER ?= go run github.com/goreleaser/goreleaser/v2@v2.18.2

.PHONY: help build install test vet fmt fmt-check check version clean snapshot release-check

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

check: fmt-check vet test ## Run formatting, vet, and race tests

version: ## Print the CLI version
	go run . --version

clean: ## Remove build artifacts
	rm -rf bin dist

snapshot: ## Build a local release snapshot with GoReleaser (no publish)
	$(GORELEASER) release --snapshot --clean --skip=publish,notarize

release-check: check ## Validate release config, archives, checksums, and local binary
	$(GORELEASER) check
	$(MAKE) snapshot
	python3 scripts/check-release.py
