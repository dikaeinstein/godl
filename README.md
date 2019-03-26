# godl

Godl is a CLI tool used to download and install go binary releases on mac.

[![Build Status](https://travis-ci.com/dikaeinstein/godl.svg?branch=master)](https://travis-ci.com/dikaeinstein/godl)

## Installation

```go get -u https://github.com/dikaeinstein/godl```

```go install github.com/dikaeinstein/godl```

If you've have setup your $GOPATH correctly, you should have `godl` command in your $PATH.

Run `godl --version` to verify.

Run `godl --help` to get help and see available options

Example: `godl go1.11.4.darwin-amd64.tar.gz ~/Downloads`

## Musings

This tool is a direct conversion of my shell script to download and install go binary releases. So its a naive implementation that gets the job done 😎.

### Improvements

The implementation of the tool could get better by taking advantage of
 the concurrency primitives of golang. Like using the `net/http` package to
 fetch the binary archive and stream the response to a goroutine in charge of extracting the tarball and installing it. Similar to this shell command: `curl -L [url] | tar -xfz [path/to/install/binary]`.
