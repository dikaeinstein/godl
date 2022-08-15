# godl

Godl is a CLI tool used to download and install go binary releases on mac. It supports both Intel and M1 macs.

[![Build Status](https://github.com/dikaeinstein/godl/actions/workflows/ci-cd.yml/badge.svg?branch=master)](https://github.com/dikaeinstein/godl/actions)
[![Coverage Status](https://coveralls.io/repos/github/dikaeinstein/godl/badge.svg?branch=master)](https://coveralls.io/github/dikaeinstein/godl?branch=master)

## Download and Install

To download and install a specific version of godl, copy and paste the installation command:

```bash
curl -s https://raw.githubusercontent.com/dikaeinstein/godl/master/get.sh | sh -s -- v0.17.0
```

## Install with Go 1.16+

From Go 1.16, the `go install` command can be used to install `godl` directly. The binary is placed in $GOPATH/bin, or in $GOBIN if set:

```bash
go install github.com/dikaeinstein/godl
```

## Build From Source

Prerequisites for building from source are:

- make
- Go 1.15+

```bash
git clone https://github.com/dikaeinstein/godl
cd godl
make fetch
make install
```

If you've have setup your $GOPATH or $GOBIN correctly, you should have `godl` command in your $PATH.

*To check the $GOPATH and $GOBIN; run `go env`*

Run `godl version` to verify.

Run `godl help` to get help and see available options

## Subcommands

*godl supports autocomplete for all subcommands*

- `completion` — Generates completion scripts for bash or zsh
- `download` — Download go binary archive
- `help` — Help about any command
- `install` — Installs the specified go binary release version
- `list|ls` — List the downloaded versions
- `list-remote|ls-remote` — List available remote versions
- `update` - Checks for updates.
- `version` — Show the godl version information

## Typical Usage (example)

To download and install go1.16:

```bash
godl download 1.16

sudo godl install 1.16
```

Then run

```bash
go version
```

```bash
output: go version go1.16 darwin/amd64 // or something similar
```

### Contributing

You can take a shot at the suggested improvements from the README. Also follow the convention from the `contribution.md`
