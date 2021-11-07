.PHONY: bench build-linux clean doc env help lint lint-fix test
SHELL			:= /bin/bash

APP_NAME ?= gobcrypt
APP_LANG ?= golang
APP_SRC ?= ./cmd

AUTHOR := @mijho
BUILD := $(shell date +%Y%m%dT%H%M%S%z)
COMMIT := $(shell git rev-parse --short HEAD)
EMAIL :=
VERSION := 2.0.0
GOLANG_VERSION ?= 1.17

LDFLAGS := -ldflags \
	'-X "main.Author=$(AUTHOR)" \
	-X main.Build=$(BUILD) \
	-X main.Commit=$(COMMIT) \
	-X main.Email=$(EMAIL) \
	-X main.Name=$(APP_NAME) \
	-X main.Version=$(VERSION)'

APP_BUILD_FLAGS := $(LDFLAGS)

## Build binary
build:
	@go build $(LDFLAGS) -v -x -o $(APP_NAME) $(APP_SRC)

## Install all required tools for development
env:
	$(info Install all required tools for development)
	$(GO_MOD)
	go get ./...

## Build for Darwin amd64 platform
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 \
go build $(LDFLAGS) -o $(APP_NAME)_darwin_amd64 $(APP_SRC)	

## Build for GNU/Linux amd64 platform
build-linux-amd64:
	GOOS=linux GOARCH=amd64 \
go build $(LDFLAGS) -o $(APP_NAME)_linux_amd64 $(APP_SRC)

## Build for GNU/Linux arm64 platform
build-linux-arm64:
	GOOS=linux GOARCH=arm64 \
go build $(LDFLAGS) -o $(APP_NAME)_linux_arm64 $(APP_SRC)

## Build and run application
run:
	@go run $(APP_SRC)

## Run available tests
test:
	@go test -v ./...

## Run static check, check code formatting and type check
lint:
	$(info Run static check, check code formating and type check)
	@golangci-lint run --verbose

## Apply linting
lint-fix:
	$(info Apply linting)
	@golangci-lint run --verbose --fix

## Extract documentation available on http://localhost:6060
doc:
	$(info Extract documenation available on http://localhost:6060)
	@godoc -http=localhost:6060

## Clean up built artifacts
clean:
	@$(RM) $(APP_NAME) $(APP_NAME)_*
	@go clean