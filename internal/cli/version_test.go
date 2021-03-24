package cli

import (
	"fmt"
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestVersionCmd(t *testing.T) {
	v := VersionOption{
		BuildDate:   "2021-03-14 00:28",
		GitHash:     "02cb593",
		GodlVersion: "v0.11.6",
		GoVersion:   "go version go1.16.2 darwin/amd64",
	}

	godl := NewRootCmd()
	version := NewVersionCmd(v)
	godl.RegisterSubCommands([]*cobra.Command{version})

	expectedOutput := fmt.Sprintf(`Version: %s
Go version: %s
Git hash: %s
Built: %s
`, v.GodlVersion, v.GoVersion, v.GitHash, v.BuildDate)

	output, errOutput := test.ExecuteCommand(t, false, godl.CobraCmd, "version")
	if errOutput != "" {
		t.Errorf("godl version failed: expected errOutput %s; got %s", "", errOutput)
	}

	if output != expectedOutput {
		t.Errorf("godl version failed: expected output %s; got %s", expectedOutput, output)
	}
}
