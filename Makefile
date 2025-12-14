# ---------------------------------------------------------------------------
# Makefile for CLI utilities
# 
# Escape '#' and '[' characters with '\', and '$' characters with '$$'
# ---------------------------------------------------------------------------

PROJECT_NAME=$(shell basename $$(pwd))
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0-dev")
BUILD_DATE=$(shell date -u "+%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags=all="-s -w -X \"main.name=$(PROJECT_NAME)\" -X \"main.version=$(VERSION)\" -X \"main.date=$(BUILD_DATE)\" -X \"main.commit=$(GIT_COMMIT)\""

# Extract version components for goversioninfo
XYZ_VERSION=$(shell echo $(VERSION) | sed -E 's/^v([0-9]+\.[0-9]+\.[0-9]+).*/\1/')
VER_MAJOR = $(shell echo $(XYZ_VERSION) | cut -d. -f1)
VER_MINOR = $(shell echo $(XYZ_VERSION) | cut -d. -f2)
VER_PATCH = $(shell echo $(XYZ_VERSION) | cut -d. -f3)
VER_BUILD = 0  # Set default build number if needed

MAKEFLAGS += --no-print-directory

## This is faster than extracting variables from git each time
FASTBUILD_CMD=go build -ldflags='-s -w -X main.name=$(PROJECT_NAME) -X main.version=v1.0.0-dirty -X main.date=$(BUILD_DATE) -X main.commit=12345'

all: dist

## build: build project
build: goversioninfo
	go build $(LDFLAGS)

## watch: watch for changes and rebuild (a first build is necessary to create .syso files)
watch: build
	@start watchexec.exe --postpone --timings --exts go \""$(FASTBUILD_CMD)"\" &

goversioninfo:
	@goversioninfo -product-name $(PROJECT_NAME) \
                  -product-version $(VERSION) \
                  -ver-major $(VER_MAJOR) \
                  -ver-minor $(VER_MINOR) \
                  -ver-patch $(VER_PATCH) \
                  -ver-build $(VER_BUILD) \
                  -o resource_386.syso
	@goversioninfo -product-name $(PROJECT_NAME) \
                  -product-version $(VERSION) \
                  -ver-major $(VER_MAJOR) \
                  -ver-minor $(VER_MINOR) \
                  -ver-patch $(VER_PATCH) \
                  -ver-build $(VER_BUILD) \
                  -64 \
                  -o resource_amd64.syso
	@goversioninfo -product-name $(PROJECT_NAME) \
                  -product-version $(VERSION) \
                  -ver-major $(VER_MAJOR) \
                  -ver-minor $(VER_MINOR) \
                  -ver-patch $(VER_PATCH) \
                  -ver-build $(VER_BUILD) \
                  -arm \
                  -o resource_arm64.syso

## test: run tests with coverage
test:
	go test -v -cover ./...

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
release: goversioninfo
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
	rm -f coverage.out
	rm -f releaseinfo
	rm -f resource_*.syso

## version: show version info
version:
	@echo "$(PROJECT_NAME) $(VERSION), built on $(BUILD_DATE) (commit: $(GIT_COMMIT))"
	@echo
	@echo "LDFLAGS:"
	@echo "    $(LDFLAGS)"
	@echo
	@echo "CHANGELOG:"
	@sed '1,/\#\# \[${VERSION}/d;/^\#\# /Q' CHANGELOG.md

## help: display this help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: all test clean help
