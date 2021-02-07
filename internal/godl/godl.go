package godl

import (
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// New returns the base command when called without any subcommands
func New() *GodlCmd {
	cobra.OnInitialize(initConfig)

	return &GodlCmd{
		cobraCmd: &cobra.Command{
			Use:   "godl [command]",
			Short: "Godl is a CLI tool used to download and install go binary releases on mac.",
		},
	}
}

type GodlCmd struct {
	cobraCmd *cobra.Command
}

// Execute adds all child commands to the base `godl` command and sets flags appropriately.
// This is called by main. It only needs to happen once.
func (godl *GodlCmd) Execute() {
	if err := godl.cobraCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (godl *GodlCmd) ExecuteC() (*cobra.Command, error) {
	return godl.cobraCmd.ExecuteC()
}

func (godl *GodlCmd) GenerateBashCompletionFile(bashTarget string) error {
	return godl.cobraCmd.GenBashCompletionFile(bashTarget)
}

func (godl *GodlCmd) GenerateZshCompletionFile(zshTarget string) error {
	return godl.cobraCmd.GenZshCompletionFile(zshTarget)
}

func (godl *GodlCmd) RegisterSubCommands(subCommands []*cobra.Command) {
	for _, subCmd := range subCommands {
		godl.cobraCmd.AddCommand(subCmd)
	}
}

func (godl *GodlCmd) SetArgs(args []string) {
	godl.cobraCmd.SetArgs(args)
}

func (godl *GodlCmd) SetErr(newErr io.Writer) {
	godl.cobraCmd.SetErr(newErr)
}

func (godl *GodlCmd) SetOut(output io.Writer) {
	godl.cobraCmd.SetOut(output)
}

func (godl *GodlCmd) SetOutput(output io.Writer) {
	godl.cobraCmd.SetOutput(output)
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
