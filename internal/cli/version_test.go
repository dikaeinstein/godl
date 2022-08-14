package cli

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/dikaeinstein/godl/test"
)

func TestVersionCmd(t *testing.T) {
	v := VersionOption{
		BuildDate:   "2021-03-14 00:28",
		GitHash:     "02cb593",
		GodlVersion: "v0.11.6",
		GoVersion:   "go1.16.2",
	}

	godl := newRootCmd()
	version := newVersionCmd(v)
	registerSubCommands(godl, []*cobra.Command{version})

	expectedOutput := fmt.Sprintf(`Version: %s
Go version: %s
Git hash: %s
Built: %s
`, v.GodlVersion, v.GoVersion, v.GitHash, v.BuildDate)

	output, errOutput := test.ExecuteCommand(t, false, godl, "version")
	require.Equal(t, "", errOutput)
	require.Equal(t, expectedOutput, output)
}
