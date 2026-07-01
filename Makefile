.PHONY: all templ build test lint vet fmt run-demo generate tidy clean golden cover ci

TEMPL_PKG := github.com/a-h/templ/cmd/templ
TEMPL_VERSION := v0.3.1020

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

## Regenerate golden snapshots (SVG + path strings) after an intentional
## render change. Run this when you know the output should change; commit the
## resulting testdata/golden/*.txt updates alongside the source change.
##
## Only packages that import internal/golden honor the -update flag; passing
## it to other packages makes `go test` reject the unknown flag, so we scope
## the regeneration to the packages that own golden snapshots.
golden:
	go test ./charts/arcs ./charts/bar ./charts/line ./charts/pie \
		./charts/heatmap ./charts/waffle ./charts/calendar \
		./charts/radar ./charts/radialbar \
		./charts/scatterplot ./charts/stream ./charts/bullet ./charts/funnel ./charts/boxplot \
		./charts/bump ./charts/marimekko ./charts/parallelcoordinates ./charts/polarbar \
		./charts/treemap ./charts/sunburst ./charts/icicle ./charts/circlepacking ./charts/tree \
		./charts/voronoi \
		-update

## Run tests with coverage, writing a coverage profile + HTML report.
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@ echo "coverage: coverage.out (text) + coverage.html (html)"

## CI entry point: lint + test (golden snapshots compared, not regenerated).
ci: lint test
	@ echo "ci: ok"

## Run the demo app
run-demo:
	go run ./examples/app

## Tidy modules
tidy:
	go mod tidy

## Clean generated templ files
clean:
	find . -type f -name 'templ_*.go' -not -path './contrib/*' -not -path './.opencode/*' -delete 2>/dev/null || true