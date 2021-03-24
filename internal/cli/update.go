package cli

import (
	"context"
	"net/http"
	"time"

	"github.com/dikaeinstein/godl/internal/godl/update"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(client *http.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Check for updates.",
		Long: `If you run into 403 Forbidden errors from Github release API,
 you need to create a config file at ~/.godl/config.json.
 Paste this into the file and replace 'yourGithubPersonalAccessToken' with your own token:
{
    "gh_token": "yourGithubPersonalAccessToken"
}
 `,
		RunE: func(cmd *cobra.Command, args []string) error {
			u := update.Update{Client: client, Output: cmd.OutOrStdout()}

			const timeout = 15
			ctx, cancel := context.WithTimeout(cmd.Context(), timeout*time.Second)
			defer cancel()

			return u.Run(ctx, godlVersion)
		},
	}
}
