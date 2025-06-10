/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// configdelCmd represents the configdel command
var configDelCmd = &cobra.Command{
	Use:     "del",
	Aliases: []string{"rm"},
	Short:   "Delete nacos server config",
	RunE: func(cmd *cobra.Command, args []string) error {
		cliConfig.DeleteServer(args[0])
		return cliConfig.WriteFile(cmdOpts.ConfigFile)
	},
}

func init() {
	configCmd.AddCommand(configDelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configdelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configdelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
