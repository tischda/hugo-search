# ----------------------------------------------------------------------------
# Makefile for hugo-search (windows specific)
# 
# Compiler: GO 1.7.3
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

# notice the absurd --no-preserve-root option needed for MSYS to recognize slashes!
# and rm.exe must be in a path with no spaces, eg. c:/sbin/git-sdk-64/usr/bin/rm.exe
# strange that rd /s/q does not work... looks like make is running a bash interpreter
clean:
	go clean
	rm --no-preserve-root -rf test/indexes
	rm --no-preserve-root -rf test/public

start:
	start hugo -s test server --port=$(HUGO_PORT)
	start hugo-search --addr=:$(BLEVE_PORT) --hugoPath=test --indexPath=test/indexes/search.bleve
	start "http://localhost:$(HUGO_PORT)/"

# mind the double slashes, this is run under /bin/sh in Windows...
stop: 
	taskkill //F //IM hugo-search.exe
	taskkill //F //IM hugo.exe

restart: stop start
