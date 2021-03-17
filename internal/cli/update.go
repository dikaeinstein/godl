package cli

import (
	"context"
	"net/http"
	"time"

	"github.com/dikaeinstein/godl/internal/godl/update"
	"github.com/spf13/cobra"
)

func NewUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Check for updates.",
		RunE: func(cmd *cobra.Command, args []string) error {
			u := update.Update{Client: http.DefaultClient, Output: cmd.OutOrStdout()}

			const timeout = 15
			ctx, cancel := context.WithTimeout(cmd.Context(), timeout*time.Second)
			defer cancel()

			return u.Run(ctx, godlVersion)
		},
	}
}
