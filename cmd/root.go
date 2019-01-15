// Copyright © 2019 Onyedikachi Solomon Okwa <solozyokwa@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "godl [go_archive] [path_to_save_archive]",
	Short: "Godl is a CLI tool used to download and install go binary releases on mac",
	Long: `
Godl is a CLI tool used to download and install go binary releases on mac.
It downloads the go binary archive specified from https://golang.org/dl/, saves it at specified
path and unpacks it into /usr/local/. The downloaded archive can be found at specified download path
or $HOME/Downloads by default.`,
	Example: "godl go1.11.4.darwin-amd64.tar.gz ~/Downloads -r",
	Version: "0.0.1",
	Run:     downloadAndInstallGo,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("provide binary archive to download")
		}
		if !strings.HasSuffix(args[0], ".tar.gz") {
			return errors.New("provide valid archive name i.e a tar.gz file")
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.godl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("remove", "r", false, "Remove flag is optional and is used to remove the downloaded archive after installing go.")
}

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

// Download and install go binary
func downloadAndInstallGo(cmd *cobra.Command, args []string) {
	var downloadPath string
	var archiveName = args[0]
	if len(args) > 1 && len(args) <= 2 {
		downloadPath = args[1] + archiveName
	} else {
		home, _ := homedir.Dir()
		downloadPath = path.Join(home, "Downloads", archiveName)
	}

	downloadGoBinary(archiveName, downloadPath)
	installGo(downloadPath)
}

func downloadGoBinary(archiveName, downloadPath string) {
	const HOST = "https://dl.google.com/go/"
	cm := exec.Command("curl", "-L", HOST+archiveName, "-o", downloadPath)
	handleShellCommand(cm)
}

func installGo(archivePath string) {
	cm := exec.Command("tar", "-C", path.Join("/usr", "local", "test"), "-xzf", archivePath)
	handleShellCommand(cm)
}

func displayProgress(rc io.ReadCloser) {
	input := bufio.NewScanner(rc)
	input.Split(bufio.ScanBytes)
	for input.Scan() {
		fmt.Print(input.Text())
	}
}

func handleShellCommand(c *exec.Cmd) {
	stderr, _ := c.StderrPipe()
	err := c.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	go displayProgress(stderr)
	err = c.Wait()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
