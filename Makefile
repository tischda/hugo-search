# ---------------------------------------------------------------------------
# Makefile for CLI utilities
# 
# Escape '#' and '[' characters with '\', and '$' characters with '$$'
# ---------------------------------------------------------------------------

BUILD_TAG=$(shell git describe --tags 2>/dev/null || echo unreleased)
LDFLAGS=-ldflags=all="-X main.version=${BUILD_TAG} -s -w"

HUGO_PORT  = 1313
BLEVE_PORT = 8080

all: build

build:
	go build ${LDFLAGS}

test:
	go test -v -cover

cover:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

install:
	go install ${LDFLAGS} ./...

update:
	go get -u
	go mod tidy
	# https://github.com/golang/go/issues/45161
	go mod vendor

snapshot:
	goreleaser --snapshot --skip-publish --rm-dist

release: 
	@sed '1,/\#\# \[${BUILD_TAG}/d;/^\#\# /Q' CHANGELOG.md > releaseinfo
	goreleaser release --rm-dist --release-notes=releaseinfo
	@rm -f releaseinfo

clean:
	go clean
	rm -f releaseinfo
	rm -rf dist
	rm -f coverage.out
	rm -rf test/indexes
	rm -rf test/public

start: clean build
	./hugo-search --addr=:$(BLEVE_PORT) --hugoPath=test --indexPath=test/indexes/search.bleve --verbose &
	@echo .
	@echo .
	@echo .
	./test/hugo -s test server --port=$(HUGO_PORT) &

stop:
	taskkill /F /IM hugo-search.exe
	taskkill /F /IM hugo.exe

restart: stop start
