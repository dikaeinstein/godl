package cli

import (
	"context"

	"github.com/dikaeinstein/godl/pkg/exitcode"
	"github.com/spf13/cobra"
)

func Run() int {
	// root command
	godl := NewRootCmd()
	// subcommands
	c := NewCompletionCmd(godl)
	d := NewDownloadCmd()
	i := NewInstallCmd()
	ls := NewListCmd()
	lsr := NewListRemoteCmd()
	u := NewUpdateCmd()
	v := NewVersionCmd()

	godl.RegisterSubCommands([]*cobra.Command{c, d, i, ls, lsr, u, v})

	return exitcode.Get(godl.Execute(context.Background()))
}
