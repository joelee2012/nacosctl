/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/joelee2012/nacosctl/pkg/nacos"
	"github.com/spf13/cobra"
)

// createCsCmd represents the createCs command
var createCsCmd = &cobra.Command{
	Use:   "cs [flags] name",
	Short: "Create one configuration",

	Run: func(cmd *cobra.Command, args []string) {
		client := NewNacosClient()
		createOpts.DataID = args[0]
		cobra.CheckErr(client.CreateConfig(&createOpts))
		fmt.Printf("configuration/%s created\n", createOpts.DataID)

	},
	Args: cobra.ExactArgs(1),
}

var createOpts nacos.CreateCSOpts

func init() {
	createCmd.AddCommand(createCsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCsCmd.Flags().StringVarP(&createOpts.NamespaceID, "namespace", "n", "", "namespace id")
	createCsCmd.Flags().StringVarP(&createOpts.Group, "group", "g", "DEFAULT_GROUP", "group of configuration")
	createCsCmd.Flags().StringVarP(&createOpts.Content, "content", "c", "", "content of configuration")
	createCsCmd.MarkFlagRequired("content")
	createCsCmd.Flags().StringVarP(&createOpts.Type, "type", "t", "text", "configuration type")
	createCsCmd.Flags().StringVarP(&createOpts.Description, "description", "d", "", "description of configuration")
	createCsCmd.Flags().StringVarP(&createOpts.Tags, "tags", "T", "", "tags of configuration")
	createCsCmd.Flags().StringVarP(&createOpts.Application, "application", "a", "", "application of configuration")

}
