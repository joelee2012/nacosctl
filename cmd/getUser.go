/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"slices"

	"github.com/joelee2012/nacosctl/pkg/nacos"
	"github.com/spf13/cobra"
)

// getUserCmd represents the getCs command
var getUserCmd = &cobra.Command{
	Use:     "user [name]",
	Aliases: []string{"user"},
	Short:   "Display one or many user",
	Run: func(cmd *cobra.Command, args []string) {
		getUser(args)
	},
}

func init() {
	getCmd.AddCommand(getUserCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getUserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

}

func getUser(args []string) {
	client := NewNacosClient()
	users, err := client.ListUser()
	cobra.CheckErr(err)
	if len(args) > 0 {
		var us []*nacos.User
		for _, u := range users.Items {
			if slices.Contains(args, u.Name) {
				us = append(us, u)
			}
		}
		users.Items = us
	}
	list := NewList(client.APIVersion, users.Items, NewUser)
	cobra.CheckErr(WriteFormat(list, cmdOpts.Output, os.Stdout))
}
