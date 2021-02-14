BINARY_NAME=godl
PACKAGE=github.com/dikaeinstein/godl/internal/godl
BUILD_DATE=$(shell date +%Y-%m-%d\ %H:%M)
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --tags)
GO_VERSION=$(shell go version)
GOBIN="$(shell go env GOPATH)/bin"

# setup -ldflags for go build
LD_FLAGS=-ldflags '-s -w -X "$(PACKAGE)/version.godlVersion=$(VERSION)" -X "$(PACKAGE)/version.buildDate=$(BUILD_DATE)" -X "$(PACKAGE)/version.goVersion=$(GO_VERSION)" -X "$(PACKAGE)/version.gitHash=$(GIT_COMMIT_HASH)"'

## Fetch dependencies
fetch:
	GO111MODULE=on go get -v ./...

lint:
	GO111MODULE=on golangci-lint run ./...

## Run tests
test:
	GO111MODULE=on go test -race ./...

## Run tests with coverage
test-cover:
	GO111MODULE=on go test -coverprofile=cover.out -race ./...

## Build binary
build:
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -a $(LD_FLAGS) -o godl cmd/main.go

## Simulate installing the binary to $GOBIN path using `go build`
install:
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -a $(LD_FLAGS) -o $(GOBIN)/godl cmd/main.go

## Execute binary
run:
	GO111MODULE=on go run -a $(LD_FLAGS) cmd/main.go

.PHONY: clean test
## Remove binary
clean:
	if [ -f $(BINARY_NAME) ]; then rm -f $(BINARY_NAME); fi
