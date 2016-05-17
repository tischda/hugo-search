# ----------------------------------------------------------------------------
# Makefile for hugo-search (windows specific)
# 
# Compiler: GO 1.6
# ----------------------------------------------------------------------------

PROJECT_DIR=$(notdir $(shell pwd))

HUGO_PORT  = 1313
BLEVE_PORT = 8080

all: get build

build:
	go build -ldflags "-X main.version=`git describe --tags` -s -w"

get:
	go get -v

test: fmt
	go test -v -cover

cover:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

fmt:
	go fmt

vet:
	go vet -v

install:
	go install -ldflags "-X main.version=`git describe --tags` -s -w"

dist: clean build
	upx -9 ${PROJECT_DIR}.exe

clean:
	go clean

start:  $(EXECUTABLE)
	cmd /c start hugo -s test server --port=$(HUGO_PORT)
	cmd /c start hugo-search --addr=:$(BLEVE_PORT) --hugoPath=test --indexPath=test/indexes/search.bleve
	cmd /c start http://localhost:$(HUGO_PORT)/

stop: 
	pskill.exe hugo-search
	pskill.exe hugo
