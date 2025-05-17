/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"slices"

	"github.com/spf13/cobra"
)

// getCsCmd represents the getCs command
var getCsCmd = &cobra.Command{
	Use:     "configurations [name]",
	Aliases: []string{"cs"},
	Short:   "Display one or many configurations",
	Run: func(cmd *cobra.Command, args []string) {
		GetCs(args)
	},
}

var listOpts ListCSOpts
var showAll bool

func init() {
	getCmd.AddCommand(getCsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCsCmd.Flags().StringVarP(&listOpts.Tenant, "namespace", "n", "", "namespace id")
	getCsCmd.Flags().StringVarP(&listOpts.Group, "group", "g", "", "group name")
	getCsCmd.Flags().IntVarP(&listOpts.PageNumber, "page-number", "p", 1, "page number")
	getCsCmd.Flags().IntVarP(&listOpts.PageSize, "page-size", "s", 10, "page size")
	getCsCmd.Flags().BoolVarP(&showAll, "all", "A", false, "show all configurations")

}

func GetCs(args []string) {
	naClient, err := NewNacosClient()
	if err != nil {
		log.Fatal(err)
	}

	var allCs ConfigList
	if showAll {
		nss, err := naClient.ListNamespace()
		if err != nil {
			log.Fatal(err)
		}
		for _, ns := range nss.Items {
			listOpts.Tenant = ns.Name
			listOpts.PageNumber = 1
			for {
				cs, err := naClient.ListConfig(&listOpts)
				if err != nil {
					log.Fatal(err)
				}
				allCs.PageItems = append(allCs.PageItems, cs.PageItems...)
				if cs.PagesAvailable == 0 || cs.PagesAvailable == cs.PageNumber {
					break
				}
				listOpts.PageNumber += 1
			}
		}
	} else {
		allCs, err := naClient.ListConfig(&listOpts)
		if err != nil {
			log.Fatal(err)
		}
		if len(args) > 0 {
			var items []*Config
			for _, ns := range allCs.PageItems {
				if slices.Contains(args, ns.DataID) {
					items = append(items, ns)
				}
			}
			allCs.PageItems = items
		}
	}
	PrintResources(&allCs, os.Stdout, output)
}
