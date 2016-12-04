# -------------------------------------------------------------------------------------------------------------------
# Makefile for hugo-search (windows specific)
# 
# Compiler: GO 1.7.4
# Make: http://win-builds.org/doku.php/1.5.0_packages#make_40-5_-_gnu_make_utility_to_maintain_groups_of_programs
# -------------------------------------------------------------------------------------------------------------------

# This does not work if set via environment variables
SHELL=/Windows/system32/cmd.exe
PROJECT_DIR=$(notdir $(shell pwd))

BUILD_TAG=$(shell git describe --tags)
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

# rm is an external : C:\Program Files\Git\usr\bin\rm.exe
clean:
	go clean
	rm -rf test/indexes
	rm -rf test/public

start: clean install
	cmd /c "start hugo -s test server --port=$(HUGO_PORT)"
	cmd /c "start hugo-search --addr=:$(BLEVE_PORT) --hugoPath=test --indexPath=test/indexes/search.bleve --verbose"
	cmd /c "start http://localhost:$(HUGO_PORT)/"

stop:
	taskkill /F /IM hugo-search.exe
	taskkill /F /IM hugo.exe

restart: stop start
