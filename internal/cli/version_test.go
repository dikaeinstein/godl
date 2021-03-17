package cli

import (
	"fmt"
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestVersionCmd(t *testing.T) {
	godl := NewRootCmd()
	version := NewVersionCmd()
	godl.RegisterSubCommands([]*cobra.Command{version})

	godlVersion = "v0.11.6"
	goVersion = "go version go1.16.2 darwin/amd64"
	gitHash = "02cb593"
	buildDate = "2021-03-14 00:28"

	expectedOutput := fmt.Sprintf(`Version: %s
Go version: %s
Git hash: %s
Built: %s
`, godlVersion, goVersion, gitHash, buildDate)

	output, errOutput := test.ExecuteCommand(t, false, godl.CobraCmd, "version")
	if errOutput != "" {
		t.Errorf("godl version failed: expected errOutput %s; got %s", "", errOutput)
	}

	if output != expectedOutput {
		t.Errorf("godl version failed: expected output %s; got %s", expectedOutput, output)
	}
}
