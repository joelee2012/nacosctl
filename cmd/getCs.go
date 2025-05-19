/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path"
	"slices"

	"github.com/spf13/cobra"
)

// getCsCmd represents the getCs command
var getCsCmd = &cobra.Command{
	Use:     "configurations [name]",
	Aliases: []string{"cs"},
	Short:   "Display one or many configurations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return GetCs(args)
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
	getCsCmd.Flags().StringVarP(&cmdOpts.Namespace, "namespace", "n", "", "namespace id")
	getCsCmd.Flags().StringVarP(&cmdOpts.Group, "group", "g", "", "group name")
	getCsCmd.Flags().BoolVarP(&cmdOpts.ShowAll, "all", "A", false, "show all configurations")

}

func GetCs(args []string) error {
	naClient, err := NewNacosClient()
	if err != nil {
		return err
	}

	allCs := new(ConfigList)
	if cmdOpts.ShowAll {
		allCs, err = naClient.ListAllConfig()
		if err != nil {
			return err
		}
	} else {
		cs, err := naClient.ListConfigInNs(cmdOpts.Namespace, cmdOpts.Group)
		if err != nil {
			return err
		}
		if len(args) > 0 {
			for _, c := range cs.PageItems {
				if slices.Contains(args, c.DataID) {
					allCs.PageItems = append(allCs.PageItems, c)
				}
			}
		} else {
			allCs = cs
		}

	}
	if cmdOpts.OutDir != "" {
		allCs.FixDefaultNs()
		for _, c := range allCs.PageItems {
			dir := path.Join(cmdOpts.OutDir, c.Tenant, c.Group)
			os.MkdirAll(dir, 0750)
			c.WriteFile(path.Join(dir, c.DataID))
		}
	} else {
		allCs.FixDefaultNs()
		WriteAsFormat(cmdOpts.Output, allCs)
	}
	return nil
}
