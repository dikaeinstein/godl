BINARY_NAME=godl
PACKAGE=github.com/dikaeinstein/godl/internal/cli
BUILD_DATE=$(shell date +%Y-%m-%d\ %H:%M)
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --tags)
GO_VERSION=$(shell go version)
GOBIN="$(shell go env GOPATH)/bin"

# setup -ldflags for go build
LD_FLAGS=-ldflags '-s -w -X "$(PACKAGE).godlVersion=$(VERSION)" -X "$(PACKAGE).buildDate=$(BUILD_DATE)" -X "$(PACKAGE).goVersion=$(GO_VERSION)" -X "$(PACKAGE).gitHash=$(GIT_COMMIT_HASH)"'

## Fetch dependencies
fetch:
	go mod download

lint:
	golangci-lint run -v

## Run tests
test:
	go test -race ./...

## Run tests with coverage
test-cover:
	go test -coverprofile=cover.out -race ./...

## Build binary
build:
	GOOS=darwin GOARCH=amd64 go build -a $(LD_FLAGS) -o godl cmd/main.go

## Simulate installing the binary to $GOBIN path using `go build`
install:
	GOOS=darwin GOARCH=amd64 go build -a $(LD_FLAGS) -o $(GOBIN)/godl cmd/main.go

## Execute binary
run:
	go run -a $(LD_FLAGS) cmd/main.go

.PHONY: clean test
## Remove binary
clean:
	if [ -f $(BINARY_NAME) ]; then rm -f $(BINARY_NAME); fi
