/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// getRoleCmd represents the getCs command
var getRoleCmd = &cobra.Command{
	Use:     "role [name]",
	Aliases: []string{"r"},
	Short:   "Display one or many role",
	Run: func(cmd *cobra.Command, args []string) {
		getRole(args)
	},
}

func init() {
	getCmd.AddCommand(getRoleCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getRoleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

func getRole(args []string) {
	client := NewNacosClient()
	roles, err := client.ListRole()
	cobra.CheckErr(err)
	list := NewList(client.APIVersion, roles.Items, NewRole)
	cobra.CheckErr(WriteFormat(list, cmdOpts.Output, os.Stdout))
}
