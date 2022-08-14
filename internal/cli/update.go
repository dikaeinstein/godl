package cli

import (
	"context"
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/dikaeinstein/godl/internal/app"
)

func newUpdateCmd(client *http.Client, v VersionOption) *cobra.Command {
	uCli := &updateCli{httpClient: client, vOpt: v}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Check for updates.",
		Long: heredoc.Doc(`
			Check for updates.

			If you run into 403 Forbidden errors from Github release API, you need to a GitHub access token.
			See: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token.

			You can set the token by setting the GODL_GH_TOKEN environment variable.
			Or you can set the token in the config file:
			{
				"gh_token": "yourGithubPersonalAccessToken"
			}
		`),
		RunE: uCli.run,
	}

	const defaultTimeout = 60 * time.Second
	updateCmd.Flags().DurationP(
		"timeout", "t", defaultTimeout, "Set the check update timeout.")

	return updateCmd
}

type updateCli struct {
	httpClient *http.Client
	vOpt       VersionOption
}

func (c *updateCli) run(cmd *cobra.Command, args []string) error {
	u := app.Update{Client: c.httpClient, Output: cmd.OutOrStdout()}

	const defaultTimeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(cmd.Context(), defaultTimeout)
	defer cancel()

	return u.Run(ctx, c.vOpt.GodlVersion)
}
