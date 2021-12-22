BINARY_NAME=godl
PACKAGE=main
BUILD_DATE=$(shell date +%Y-%m-%d\ %H:%M)
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --tags)
GO_VERSION=$(shell go env GOVERSION)

# setup -ldflags for go build
LDFLAGS=-ldflags '-s -w -X "$(PACKAGE).godlVersion=$(VERSION)" -X "$(PACKAGE).buildDate=$(BUILD_DATE)" -X "$(PACKAGE).goVersion=$(GO_VERSION)" -X "$(PACKAGE).gitHash=$(GIT_COMMIT_HASH)"'

## Fetch dependencies
fetch:
	@echo Download go.mod dependencies
	@go mod download

lint:
	@golangci-lint run

## Run tests
test:
	@gotest -race $(TESTFLAGS) ./...

## Run tests with coverage
test-cover:
	@gotest -coverprofile=cover.out -race $(TESTFLAGS) ./...

## Build binary
build:
	@GOOS=darwin GOARCH=amd64 go build -a $(LDFLAGS) -o godl cmd/main.go

## Simulate installing the binary to $GOBIN path using `go build`
install:
	@GOOS=darwin GOARCH=amd64 go build -a $(LDFLAGS) -o $(shell go env GOVERSION)/godl cmd/main.go

install-tools: fetch
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

## Execute binary
run:
	@go run -a $(LDFLAGS) cmd/main.go

.PHONY: build clean fetch install install-tools lint run test test-cover

## Remove binary
clean:
	if [ -f $(BINARY_NAME) ]; then rm -f $(BINARY_NAME); fi
