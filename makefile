BINARY_NAME=godl
PACKAGE=main
BUILD_DATE=$(shell date +%Y-%m-%d\ %H:%M)
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --tags)
GO_VERSION=$(shell go env GOVERSION)

# setup -ldflags for go build
LDFLAGS=-ldflags '-s -w \
	-X "$(PACKAGE).godlVersion=$(VERSION)" \
	-X "$(PACKAGE).buildDate=$(BUILD_DATE)" \
	-X "$(PACKAGE).goVersion=$(GO_VERSION)" \
	-X "$(PACKAGE).gitHash=$(GIT_COMMIT_HASH)"'

lint:
	@golangci-lint run

## Run tests
test:
	@go run github.com/rakyll/gotest -race $(TESTFLAGS) ./...

## Run tests with coverage
test-cover:
	@go run github.com/rakyll/gotest -coverprofile=cover.out -race $(TESTFLAGS) ./...

## send test coverage to coveralls
coveralls:
	@go run github.com/mattn/goveralls -coverprofile=cover.out -service=github

## Build binary
build:
	@go build -a $(LDFLAGS) ./cmd/godl

## Build the binary to $GOBIN path using `go build`
build-install:
	@go build -a $(LDFLAGS) -o $(shell go env GOBIN)/godl cmd/main.go

## installing the binary to $GOBIN using `go install`
install:
	@go install $(LDFLAGS) ./cmd/godl

install-tools: fetch
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

## Execute binary
run:
	@go run -a $(LDFLAGS) ./cmd/godl

.PHONY: build clean fetch install install-tools lint run test test-cover

## Remove binary
clean:
	if [ -f $(BINARY_NAME) ]; then rm -f $(BINARY_NAME); fi
