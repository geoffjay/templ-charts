.PHONY: all templ build test lint vet fmt run-demo generate tidy clean

TEMPL_PKG := github.com/a-h/templ/cmd/templ
TEMPL_VERSION := v0.3.857

# Find all directories containing .templ files
TEMPL_SOURCES := $(shell find . -type f -name '*.templ' -not -path './contrib/*' -not -path './.opencode/*' 2>/dev/null)

# Go module dirs (exclude contrib, .opencode, examples/app module)
GO_PKGS := $(shell go list ./... 2>/dev/null | grep -v '/contrib/' | grep -v '/.opencode/')

all: build

## Install templ CLI if not present
$(HOME)/go/bin/templ:
	go install $(TEMPL_PKG)@$(TEMPL_VERSION)

## Generate Go code from .templ files
templ: $(HOME)/go/bin/templ
	@ if [ -n "$(TEMPL_SOURCES)" ]; then templ generate; else echo "no .templ files yet"; fi

generate: templ

## Build all Go packages
build:
	go build ./...

## Run tests
test:
	go test ./...

## Run go vet
vet:
	go vet ./...

## Run gofmt (check only)
fmt:
	@ out=$$(gofmt -l . 2>/dev/null | grep -v '/contrib/' | grep -v '/.opencode/' | grep -v '/.git/'); if [ -n "$$out" ]; then echo "gofmt needs to format:"; echo "$$out"; exit 1; else echo "gofmt: ok"; fi

## Run all linters
lint: vet fmt

## Run the demo app
run-demo:
	go run ./examples/app

## Tidy modules
tidy:
	go mod tidy

## Clean generated templ files
clean:
	find . -type f -name 'templ_*.go' -not -path './contrib/*' -not -path './.opencode/*' -delete 2>/dev/null || true