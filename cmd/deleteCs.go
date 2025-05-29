/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCsCmd represents the deleteCs command
var deleteCsCmd = &cobra.Command{
	Use:     "configurations",
	Aliases: []string{"cs"},
	Short:   "Delete one or many configurations",
	Run: func(cmd *cobra.Command, args []string) {
		naClient, err := NewNacosClient()
		cobra.CheckErr(err)
		for _, dataId := range args {
			err := naClient.DeleteConfig(&DeleteCSOpts{
				DataID: dataId,
				Group:  cmdOpts.Group,
				Tenant: cmdOpts.Namespace,
			})
			cobra.CheckErr(err)
			fmt.Printf("configuration/%s deleted\n", dataId)

		}
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	deleteCmd.AddCommand(deleteCsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deleteCsCmd.Flags().StringVarP(&cmdOpts.Namespace, "namespace", "n", "", "namespace id")
	deleteCsCmd.Flags().StringVarP(&cmdOpts.Group, "group", "g", "", "name of group")
	deleteCsCmd.MarkFlagRequired("group")
}
