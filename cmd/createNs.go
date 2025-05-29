/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createNsCmd represents the createNs command
var createNsCmd = &cobra.Command{
	Use:   "ns name",
	Short: "Create one namespace",
	Run: func(cmd *cobra.Command, args []string) {
		client := NewNacosClient()
		nsOpts.Name = args[0]
		cobra.CheckErr(client.CreateNamespace(&nsOpts))
		fmt.Printf("namespace/%#v created\n", nsOpts)
	},
	Args: cobra.ExactArgs(1),
}

var nsOpts CreateNSOpts

func init() {
	createCmd.AddCommand(createNsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createNsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createNsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// createNsCmd.Flags().StringVarP(&nsOpts.Name, "name", "n", "", "name of namespace")
	createNsCmd.Flags().StringVarP(&nsOpts.ID, "id", "i", "", "id of namespace")
	createNsCmd.MarkFlagRequired("id")
	createNsCmd.Flags().StringVarP(&nsOpts.Desc, "desc", "d", "", "description of namespace")
	createNsCmd.MarkFlagRequired("desc")
}
