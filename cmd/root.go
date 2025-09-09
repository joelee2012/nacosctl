/*
Copyright © 2025 Joe Lee <lj_2005@163.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joelee2012/nacosctl/pkg/nacos"
	"github.com/spf13/cobra"
)

type CmdOpts struct {
	NamespaceID string
	Group       string
	Output      string
	OutDir      string
	ConfigFile  string
	ShowAll     bool
}

var cmdOpts = CmdOpts{}
var cliConfig = CLIConfig{}

// var client *nacos.nacos

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nctl [options]",
	Short: "Command line tools for Nacos",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cmdOpts.ConfigFile, "setting", "s", "", "config file (default is $HOME/.nacos.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cmdOpts.ConfigFile == "" {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cmdOpts.ConfigFile = filepath.Join(home, ".nacos.yaml")
	}

	err := cliConfig.ReadFile(cmdOpts.ConfigFile)
	cobra.CheckErr(err)
}

func NewNacosClient() *nacos.Client {
	if cliConfig.Context == "" {
		cobra.CheckErr(fmt.Errorf("no context set in config file: %s", cmdOpts.ConfigFile))
	}
	server := cliConfig.GetCurrentServer()
	client := nacos.NewClient(server.URL, server.User, server.Password)
	client.DetectAPIVersion()
	return client
}
