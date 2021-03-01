package cli

import (
	"fmt"
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

func (godl *RootCmd) Execute() error {
	return godl.CobraCmd.Execute()
}

func (godl *RootCmd) GenerateBashCompletionFile(bashTarget string) error {
	return godl.CobraCmd.GenBashCompletionFile(bashTarget)
}

func (godl *RootCmd) GenerateZshCompletionFile(zshTarget string) error {
	return godl.CobraCmd.GenZshCompletionFile(zshTarget)
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
