/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configViewCmd represents the configView command
var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View nacos config",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := cliConfig.ToYaml()
		if err != nil {
			return err
		}
		_, err = fmt.Println(string(data))
		return err
	},
}

func init() {
	configCmd.AddCommand(configViewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configViewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configViewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
