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

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## lint: lint the project
lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run

## test: run tests
test:
	@go run github.com/rakyll/gotest -race $(TESTFLAGS) ./...

## test/cover: run tests with coverage
test/cover:
	@go run github.com/rakyll/gotest -coverprofile=coverage.out -race $(TESTFLAGS) ./...
	@go tool cover -html=coverage.out -o coverage.html

## coveralls: send test coverage to coveralls
coveralls:
	@go run github.com/mattn/goveralls -coverprofile=coverage.out -service=github

## build: build the binary
build:
	@go build $(LDFLAGS) ./cmd/$(BINARY_NAME)

## build/install: build the binary and outputs it tothe $GOBIN path using `go build`
build/install:
	@go build $(LDFLAGS) -o $(shell go env GOBIN)/godl ./cmd/$(BINARY_NAME)

## install: install the binary to $GOBIN using `go install`
install:
	@go install $(LDFLAGS) ./cmd/$(BINARY_NAME)

## run: execute binary
run:
	@go run $(LDFLAGS) ./cmd/$(BINARY_NAME)

.PHONY: build clean fetch install install-tools lint run test test-cover

## clean: remove binary
clean:
	if [ -f $(BINARY_NAME) ]; then rm -f $(BINARY_NAME); fi
