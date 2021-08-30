# generated-from:0f1cfb3f9faa0c83355794c5720cb80c30b77f4fcb2887d31d2887bd169db413 DO NOT REMOVE, DO UPDATE

PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
PWD := $(shell pwd)

ifndef VERSION
	VERSION := $(shell git describe --tags --abbrev=0)
endif

COMMIT_HASH :=$(shell git rev-parse --short HEAD)
DEV_VERSION := dev-${COMMIT_HASH}

USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

export GOPRIVATE=github.com/moov-io

all: install update build

.PHONY: install
install:
	go mod tidy
	go mod vendor

build:
	go build -mod=vendor -ldflags "-X github.com/moov-io/iso8583.Version=${VERSION}" -o bin/iso8583 github.com/moov-io/iso8583/cmd/iso8583

.PHONY: setup
setup:
	docker-compose up -d --force-recreate --remove-orphans

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	GOLANGCI_LINTERS=gosec COVER_THRESHOLD=75.0 ./lint-project.sh
endif

.PHONY: teardown
teardown:
	-docker-compose down --remove-orphans

docker: docker-fuzz

docker-fuzz:
	docker build --pull -t moov/iso8583fuzz:$(VERSION) . -f Dockerfile-fuzz
	docker tag moov/iso8583fuzz:$(VERSION) moov/iso8583fuzz:latest

docker-push:
	docker push moov/iso8583fuzz:${VERSION}

# Extra utilities not needed for building

run: update build
	./bin/iso8583

test: update
	go test -cover github.com/moov-io/iso8583/...

.PHONY: clean
clean:
ifeq ($(OS),Windows_NT)
	@echo "Skipping cleanup on Windows, currently unsupported."
else
	@rm -rf cover.out coverage.txt misspell* staticcheck*
	@rm -rf ./bin/
endif

# For open source projects

# From https://github.com/genuinetools/img
.PHONY: AUTHORS
AUTHORS:
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@echo "$(shell git log --format='\n%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf)" >> $@

dist: clean build
ifeq ($(OS),Windows_NT)
	CGO_ENABLED=1 GOOS=windows go build -o bin/iso8583.exe ./cmd/iso8583/
else
	CGO_ENABLED=1 GOOS=$(PLATFORM) go build -o bin/iso8583-$(PLATFORM)-amd64 ./cmd/iso8583/
endif
