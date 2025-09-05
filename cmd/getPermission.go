/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// getPermCmd represents the getCs command
var getPermCmd = &cobra.Command{
	Use:     "permission [name]",
	Aliases: []string{"perm"},
	Short:   "Display one or many permission",
	Run: func(cmd *cobra.Command, args []string) {
		getPermission(args)
	},
}

func init() {
	getCmd.AddCommand(getPermCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getPermCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

}

func getPermission(args []string) {
	client := NewNacosClient()
	perms, err := client.ListPermission()
	cobra.CheckErr(err)
	list := NewList(client.APIVersion, perms.Items, NewPermission)
	cobra.CheckErr(WriteFormat(list, cmdOpts.Output, os.Stdout))
}
