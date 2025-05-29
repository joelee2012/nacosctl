/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"slices"

	"github.com/spf13/cobra"
)

// getNsCmd represents the getNs command
var getNsCmd = &cobra.Command{
	Use:     "ns [name]",
	Aliases: []string{"namespace"},
	Short:   "Display one or many namespaces",
	Run: func(cmd *cobra.Command, args []string) {
		GetNamespace(args)
	},
}

func init() {
	getCmd.AddCommand(getNsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getNsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getNsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GetNamespace(args []string) {
	client := NewNacosClient()
	nss, err := client.ListNamespace()
	cobra.CheckErr(err)
	if len(args) > 0 {
		var items []*Namespace
		for _, ns := range nss.Items {
			if slices.Contains(args, ns.ShowName) {
				items = append(items, ns)
			}
		}
		nss.Items = items
	}
	if cmdOpts.OutDir != "" {
		cobra.CheckErr(nss.WriteToDir(cmdOpts.OutDir))
	} else {
		WriteAsFormat(cmdOpts.Output, nss)
	}
}
