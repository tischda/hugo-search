# ----------------------------------------------------------------------------
# Makefile for hugo-search (windows specific)
# 
# Compiler: GO 1.7.1
# ----------------------------------------------------------------------------

PROJECT_DIR=$(notdir $(shell pwd))

BUILD_TAG=`git describe --tags 2>/dev/null`
LDFLAGS=-ldflags "-X main.version=${BUILD_TAG} -s -w"

HUGO_PORT  = 1313
BLEVE_PORT = 8080

all: get build

build:
	go build ${LDFLAGS}

get:
	go get -v

test: clean vet
	go test -v -cover

cover:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

fmt:
	go fmt

vet:
	go vet -v

install:
	go install ${LDFLAGS}

dist: clean build
	upx -9 ${PROJECT_DIR}.exe

# mind stupid option for MSYS to recognize slashes!
clean:
	go clean
	rm -rf --no-preserve-root test/indexes
	rm -rf --no-preserve-root test/public

start:
	start hugo -s test server --port=$(HUGO_PORT)
	start hugo-search --addr=:$(BLEVE_PORT) --hugoPath=test --indexPath=test/indexes/search.bleve
	start "http://localhost:$(HUGO_PORT)/"

# mind the double slashes, this is run under /bin/sh in Windows...
stop: 
	taskkill //F //IM hugo-search.exe
	taskkill //F //IM hugo.exe

restart: stop start
