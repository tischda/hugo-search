# ---------------------------------------------------------------------------
# Makefile for CLI utilities
# 
# Escape '#' and '[' characters with '\', and '$' characters with '$$'
# ---------------------------------------------------------------------------

PROJECT_NAME=$(shell git rev-parse --show-toplevel | xargs basename )
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0-dev")
BUILD_DATE=$(shell date -u "+%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags=all="-s -w -X \"main.name=$(PROJECT_NAME)\" -X \"main.version=$(VERSION)\" -X \"main.date=$(BUILD_DATE)\" -X \"main.commit=$(GIT_COMMIT)\""


MAKEFLAGS += --no-print-directory

all: build

## build: build project
build:
	go build $(LDFLAGS)

## test: run tests with coverage
test:
	go test -v -cover ./...

## watch: watch for modifications in go files and rebuild if changed
watch:
	watchexec.exe --quiet --postpone --exts go make build

## cover: run tests and show coverage report in browser
cover:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

## install: build and install binary into workspace bin folder
install:
	go install $(LDFLAGS) ./...

## update: update dependencies
update:
	go get -u
	go mod tidy
	@# 'go mod tidy' should update the vendor directory (https://github.com/golang/go/issues/45161)
	go mod vendor

## snapshot: make a snapshot release
snapshot:
	goreleaser --snapshot --skip-publish --clean

## release: make a release based on latest tag
release: 
	@echo releasing $(VERSION)
	@sed '1,/\#\# \[${VERSION}/d;/^\#\# /Q' CHANGELOG.md > releaseinfo
	@cat releaseinfo
	@echo ----
	@goreleaser release --clean --release-notes=releaseinfo
	@rm -f releaseinfo

## dist: clean and build
dist: clean build

## clean: remove temporary files
clean:
	go clean
	rm -rf dist
	rm -f releaseinfo
	rm -f coverage.out

## version: show version info
version:
	@echo "$(PROJECT_NAME) $(VERSION), built on $(BUILD_DATE) (commit: $(GIT_COMMIT))"
	@echo "LDFLAGS:"
	@echo "    $(LDFLAGS)"

## help: display this help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: all test clean help
