# godl

Godl is a CLI tool used to download and install go binary releases on mac.

[![Build Status](https://travis-ci.com/dikaeinstein/godl.svg?branch=master)](https://travis-ci.com/dikaeinstein/godl)
[![Coverage Status](https://coveralls.io/repos/github/dikaeinstein/godl/badge.svg?branch=master)](https://coveralls.io/github/dikaeinstein/godl?branch=master)

## Installation

```go get -u https://github.com/dikaeinstein/godl```

```cd path/installed```

Then run: ```make install```

If you've have setup your $GOPATH and $GOBIN correctly, you should have `godl` command in your $PATH.

Run `godl version` to verify.

Run `godl help` to get help and see available options

## Subcommands

* `download` — Download go binary archive
* `help` — Help about any command
* `install` — Installs the specified go binary archive version
* `list` — List the downloaded versions
* `version` — Show the godl version information

## Typical Usage (example)

To download and install go1.12.4:

```bash
godl download 1.12.4

godl install 1.12.4
```

Then run

```bash
go version
```

```bash
output: go version go1.12.4 darwin/amd64 // or something similar
```

### Improvements / Coming features

* List remote versions of go
* The install command downloads archive if has not being downloaded before installing instead of failing with an error

### Contributing

You can take a shot at the suggested improvements from the README. Also follow the convention from the `contribution.md`
