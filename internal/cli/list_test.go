package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/dikaeinstein/godl/test"
)

func TestListCmd(t *testing.T) {
	ls := newListCmd()
	godl := newRootCmd()
	registerSubCommands(godl, []*cobra.Command{ls})

	_, errOutput := test.ExecuteCommand(t, false, godl, "list")

	require.Equal(t, "", errOutput)
}
