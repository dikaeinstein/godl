package cli

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/MakeNowJust/heredoc"
	"github.com/dikaeinstein/godl/pkg/text"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd() *RootCmd {
	godl := &RootCmd{
		CobraCmd: &cobra.Command{
			Use:          "godl [command]",
			Short:        "Godl is a CLI tool used to download and install go binary releases on mac.",
			SilenceUsage: true,
		},
	}
	debug := godl.CobraCmd.PersistentFlags().Bool("debug", false, "Used to turn on debug mode.")
	cobra.OnInitialize(func() { initConfig(*debug) })

	godl.CobraCmd.SetUsageTemplate(usageTemplate())

	return godl
}

type RootCmd struct {
	CobraCmd *cobra.Command
}

func (godl *RootCmd) Execute(ctx context.Context) error {
	return godl.CobraCmd.ExecuteContext(ctx)
}

func (godl *RootCmd) GenerateBashCompletion(out io.Writer) error {
	return godl.CobraCmd.GenBashCompletion(out)
}

func (godl *RootCmd) GenerateFishCompletion(out io.Writer, includeDesc bool) error {
	return godl.CobraCmd.GenFishCompletion(out, includeDesc)
}

func (godl *RootCmd) GenerateZshCompletion(out io.Writer) error {
	return godl.CobraCmd.GenZshCompletion(out)
}

// RegisterSubCommands adds all child commands to the root `godl` command
func (godl *RootCmd) RegisterSubCommands(subCommands []*cobra.Command) {
	for _, subCmd := range subCommands {
		godl.CobraCmd.AddCommand(subCmd)
	}
}

var cfgFile string

// initConfig reads in config file and ENV variables if set.
func initConfig(debug bool) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search and set config `/home/.godl/config`.
		viper.AddConfigPath(path.Join(home, ".godl"))
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("godl")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && debug {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
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
