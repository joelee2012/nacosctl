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

// getCsCmd represents the getCs command
var getCsCmd = &cobra.Command{
	Use:     "cs [name]",
	Aliases: []string{"configuration"},
	Short:   "Display one or many configurations",
	Run: func(cmd *cobra.Command, args []string) {
		GetCs(args)
	},
}

func init() {
	getCmd.AddCommand(getCsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCsCmd.Flags().StringVarP(&cmdOpts.NamespaceId, "namespace", "n", "", "namespace id")
	getCsCmd.Flags().StringVarP(&cmdOpts.Group, "group", "g", "", "group name")
	getCsCmd.Flags().BoolVarP(&cmdOpts.ShowAll, "all", "A", false, "show all configurations")

}

func GetCs(args []string) {
	client := NewNacosClient()

	allCs := new(nacos.ConfigList)
	var err error
	if cmdOpts.ShowAll {
		allCs, err = client.ListAllConfig()
		cobra.CheckErr(err)
	} else {
		cs, err := client.ListConfigInNs(cmdOpts.NamespaceId, cmdOpts.Group)
		cobra.CheckErr(err)
		if len(args) > 0 {
			for _, c := range cs.Items {
				if slices.Contains(args, c.DataID) {
					allCs.Items = append(allCs.Items, c)
				}
			}
		} else {
			allCs = cs
		}
	}
	if cmdOpts.OutDir != "" {
		cobra.CheckErr(allCs.WriteToDir(cmdOpts.OutDir))
	} else {
		nacos.WriteAsFormat(cmdOpts.Output, allCs, os.Stdout)
	}
}
