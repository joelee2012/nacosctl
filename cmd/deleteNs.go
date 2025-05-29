/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteNsCmd represents the deleteNs command
var deleteNsCmd = &cobra.Command{
	Use:     "namespace",
	Aliases: []string{"ns"},
	Short:   "Delete one or many namespaces",
	Run: func(cmd *cobra.Command, args []string) {
		naClient, err := NewNacosClient()
		cobra.CheckErr(err)
		for _, ns := range args {
			cobra.CheckErr(naClient.DeleteNamespace(ns))
			fmt.Printf("namespace/%s deleted\n", ns)
		}
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	deleteCmd.AddCommand(deleteNsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteNsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteNsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
