BINARY_NAME=godl
PACKAGE=github.com/dikaeinstein/godl
BUILD_DATE=$(shell date +%Y-%m-%d\ %H:%M)
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --tags)
GO_VERSION=$(shell go version)

# setup -ldflags for go build
LD_FLAGS=-ldflags '-X "$(PACKAGE)/cmd.version=$(VERSION)" -X "$(PACKAGE)/cmd.buildDate=$(BUILD_DATE)" -X "$(PACKAGE)/cmd.goVersion=$(GO_VERSION)" -X "$(PACKAGE)/cmd.gitHash=$(GIT_COMMIT_HASH)"'

## Fetch dependencies
install:
	GO111MODULE=on go get -v ./...

## Run tests
test:
	GO111MODULE=on APP_ENV=test go test -race -v ./...

## Build binary
build:
	GO111MODULE=on go build $(LD_FLAGS)

## Execute binary
run:build
	./$(BINARY_NAME)

.PHONY: clean
## Remove binary
clean:
	if [ -f $(BINARY_NAME) ]; then rm -f $(BINARY_NAME) fi;
