package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd() *RootCmd {
	cobra.OnInitialize(initConfig)

	return &RootCmd{
		CobraCmd: &cobra.Command{
			Use:   "godl [command]",
			Short: "Godl is a CLI tool used to download and install go binary releases on mac.",
		},
	}
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
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".godl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".godl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
