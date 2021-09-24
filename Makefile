# ---------------------------------------------------------------------------
# Makefile for CLI utilities
# ---------------------------------------------------------------------------

BUILD_TAG=$(shell git describe --tags 2>/dev/null || echo undefined)
LDFLAGS=-ldflags=all="-X main.version=${BUILD_TAG} -s -w"

HUGO_PORT  = 1313
BLEVE_PORT = 8080

all: build

build:
	go build ${LDFLAGS}

test:	clean
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
	goreleaser release --rm-dist

clean:
	go clean
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
