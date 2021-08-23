PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')

ifndef VERSION
	VERSION := $(shell git describe --tags --abbrev=0)
endif

COMMIT_HASH :=$(shell git rev-parse --short HEAD)
DEV_VERSION := dev-${COMMIT_HASH}
PWD := $(shell pwd)


.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	./lint-project.sh
endif

docker: clean docker-fuzz

docker-fuzz:
	docker build --pull -t moov/iso8583fuzz:$(VERSION) . -f Dockerfile-fuzz
	docker tag moov/iso8583fuzz:$(VERSION) moov/iso8583fuzz:latest

build:
	@mkdir -p ./bin/
	go build -ldflags "-X github.com/moov-io/iso8583.Version=${VERSION}" -o bin/iso8583 github.com/moov-io/iso8583/cmd/iso8583

release-push:
	docker push moov/iso8583fuzz:$(VERSION)

.PHONY: clean
clean:
ifeq ($(OS),Windows_NT)
	@echo "Skipping cleanup on Windows, currently unsupported."
else
	@rm -rf ./bin/ coverage.txt misspell* staticcheck lint-project.sh
endif

.PHONY: cover-test cover-web
cover-test:
	go test -coverprofile=cover.out ./...
cover-web:
	go tool cover -html=cover.out

# From https://github.com/genuinetools/img
.PHONY: AUTHORS
AUTHORS:
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@echo "$(shell git log --format='\n%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf)" >> $@
