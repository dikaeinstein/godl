//go:build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/mattn/goveralls"
	_ "golang.org/x/vuln/cmd/govulncheck"
)
