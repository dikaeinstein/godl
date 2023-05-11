package cli

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dikaeinstein/godl/internal/app"
	"github.com/dikaeinstein/godl/pkg/exitcode"
	"github.com/dikaeinstein/godl/pkg/text"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "godl [command]",
		Short:        "Godl is a CLI tool used to download and install go binary releases on mac.",
		SilenceUsage: true,
		PreRunE:      setupConfig,
	}
	rootCmd.SetUsageTemplate(usageTemplate())

	setupFlags(rootCmd)

	return rootCmd
}

func setupFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().
		Bool("debug", false, "Used to turn on debug mode.")
	cmd.PersistentFlags().
		String("config-file", "", "config file (default is $HOME/.godl/config)")
}

func registerSubCommands(root *cobra.Command, subCmds []*cobra.Command) {
	for _, subCmd := range subCmds {
		root.AddCommand(subCmd)
	}
}

func setupConfig(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return err
	}

	cfgFile := viper.GetString("config-file")

	if cfgFile != "" {
		// use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// use default config file path
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		defaultCfgFile := filepath.Join(home, ".godl", "config.json")
		viper.SetConfigFile(defaultCfgFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		// it's ok if config file doesn't exist
		e := &viper.ConfigFileNotFoundError{}
		if !errors.As(err, e) {
			return err
		}
	}

	viper.SetEnvPrefix("godl")
	viper.AutomaticEnv()

	return nil
}

func usageTemplate() string {
	return heredoc.Docf(`%s:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

%s:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

%s:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} help [command]" or "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`, text.Bold("USAGE"), text.Bold("ALIASES"), text.Bold("EXAMPLES"),
		text.Bold("AVAILABLE COMMANDS"), text.Bold("FLAGS"),
		text.Bold("INHERITED FLAGS"), text.Bold("ADDITIONAL HELP TOPICS"))
}

func Run(info app.BuildInfo) int {
	godl := newRootCmd()

	// subcommands
	completionCmd := newCompletionCmd()
	downloadCmd := newDownloadCmd(http.DefaultClient)
	installCmd := newInstallCmd(http.DefaultClient)
	lsCmd := newListCmd()
	lsRemoteCmd := newListRemoteCmd(http.DefaultClient)
	updateCmd := newUpdateCmd(http.DefaultClient, info)
	versionCmd := newVersionCmd(info)

	registerSubCommands(godl, []*cobra.Command{
		completionCmd,
		downloadCmd,
		installCmd,
		lsRemoteCmd,
		lsCmd,
		updateCmd,
		versionCmd,
	})

	return exitcode.Get(godl.ExecuteContext(context.Background()))
}
