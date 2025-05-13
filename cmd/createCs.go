/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCsCmd represents the createCs command
var createCsCmd = &cobra.Command{
	Use:   "cs [flags] name",
	Short: "Create one configuration",

	Run: func(cmd *cobra.Command, args []string) {
		naClient, err := NewNacosClient()
		if err != nil {
			fmt.Println(err)
			return
		}
		createOpts.DataID = args[0]
		if err := naClient.CreateConfig(&createOpts); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("configuration/%s created\n", createOpts.DataID)
		}
	},
	Args: cobra.ExactArgs(1),
}

var createOpts CreateCSOpts

func init() {
	createCmd.AddCommand(createCsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCsCmd.Flags().StringVarP(&createOpts.Tenant, "namespace", "n", "", "namespace id")
	createCsCmd.Flags().StringVarP(&createOpts.Group, "group", "g", "DEFAULT_GROUP", "group of configuration")
	createCsCmd.Flags().StringVarP(&createOpts.Content, "content", "c", "", "content of configuration")
	createCsCmd.MarkFlagRequired("content")
	createCsCmd.Flags().StringVarP(&createOpts.Type, "type", "t", "text", "configuration type")
}
