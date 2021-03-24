package cli

import (
	"context"
	"net/http"

	"github.com/dikaeinstein/godl/pkg/exitcode"
	"github.com/spf13/cobra"
)

func Run() int {
	// root command
	godl := NewRootCmd()
	// subcommands
	c := NewCompletionCmd(godl)
	d := NewDownloadCmd(http.DefaultClient)
	i := NewInstallCmd(http.DefaultClient)
	ls := NewListCmd()
	lsr := NewListRemoteCmd(http.DefaultClient)
	u := NewUpdateCmd(http.DefaultClient)
	v := NewVersionCmd()

	godl.RegisterSubCommands([]*cobra.Command{c, d, i, ls, lsr, u, v})

	return exitcode.Get(godl.Execute(context.Background()))
}
