# godl

Godl is a CLI tool used to download and install go binary releases on mac.

[![Build Status](https://travis-ci.com/dikaeinstein/godl.svg?branch=master)](https://travis-ci.com/dikaeinstein/godl)
[![Coverage Status](https://coveralls.io/repos/github/dikaeinstein/godl/badge.svg?branch=master)](https://coveralls.io/github/dikaeinstein/godl?branch=master)

## Standalone

godl can be easily installed as an executable. Download the latest [compiled binaries](https://github.com/dikaeinstein/godl/releases) and put it anywhere in your executable path.

## Build From Source

Prerequisites for building from source are:

- make
- Go 1.10+

```git clone https://github.com/dikaeinstein/godl```

```cd godl```

Then run: ```make install```

If you've have setup your $GOPATH and $GOBIN correctly, you should have `godl` command in your $PATH.

Run `godl version` to verify.

Run `godl help` to get help and see available options

## Subcommands

*godl supports autocomplete for all subcommands*

* `completion` — Generates completion scripts for bash or zsh
* `download` — Download go binary archive
* `help` — Help about any command
* `install` — Installs the specified go binary release version
* `list|ls` — List the downloaded versions
* `list-remote|ls-remote` — List available remote versions
* `version` — Show the godl version information

## Typical Usage (example)

To download and install go1.13:

```bash
godl download 1.13

godl install 1.13
```

Then run

```bash
go version
```

```bash
output: go version go1.13 darwin/amd64 // or something similar
```

### Improvements / Coming features

### Contributing

You can take a shot at the suggested improvements from the README. Also follow the convention from the `contribution.md`
