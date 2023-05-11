# godl

Godl is a CLI tool used to download and install go binary releases on mac. It supports both Intel and M1 macs.

[![Build Status](https://github.com/dikaeinstein/godl/actions/workflows/ci-cd.yml/badge.svg?branch=master)](https://github.com/dikaeinstein/godl/actions)
[![Coverage Status](https://coveralls.io/repos/github/dikaeinstein/godl/badge.svg?branch=master)](https://coveralls.io/github/dikaeinstein/godl?branch=master)

## Download and Install

To download and install a specific version of godl, copy and paste the installation command:

```bash
curl -s https://raw.githubusercontent.com/dikaeinstein/godl/master/get.sh | sh -s -- v0.18.0
```

## Install with Go

```bash
go install github.com/dikaeinstein/godl
```

## Build From Source

Prerequisites for building from source are:

- make
- Go 1.18+

```bash
git clone https://github.com/dikaeinstein/godl
cd godl
make build
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

To download and install go1.20:

```bash
 # To download only
godl download 1.20

 # To download and install if the given version is already downloaded.
sudo godl install 1.20

```

Then run

```bash
go version
```

```bash
output: go version go1.20 darwin/amd64 // or something similar
```

### Improvement

- update sub command to install latest version or specific version (tip: exec get.sh)
- check-update subcommand to check for an update
- support M1 Macs

### Contributing

You can take a shot at the suggested improvements from the README. Also follow the convention from the `contribution.md`
