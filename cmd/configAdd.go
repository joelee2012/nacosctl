/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// configAddCmd represents the configAdd command
var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new nacos server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cliConfig.AddServer(args[0], server)
		return cliConfig.WriteFile(cmdOpts.ConfigFile)
	},
	Args: cobra.ExactArgs(1),
}

var server = &Server{}

func init() {
	configCmd.AddCommand(configAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	configAddCmd.Flags().StringVar(&server.URL, "url", "", "the nacos url")
	configAddCmd.MarkFlagRequired("url")
	configAddCmd.Flags().StringVarP(&server.User, "user", "u", "", "nacos user")
	configAddCmd.MarkFlagRequired("user")
	configAddCmd.Flags().StringVarP(&server.Password, "password", "p", "", "nacos password")
	configAddCmd.MarkFlagRequired("password")
}
