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
	"github.com/spf13/cobra"
)

// configUseCmd represents the configSet command
var configUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Change config context",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cliConfig.SetContext(args[0]); err != nil {
			return err
		}
		return cliConfig.WriteFile(cmdOpts.ConfigFile)
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	configCmd.AddCommand(configUseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configUseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configUseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
