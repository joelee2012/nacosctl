/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/ghodss/yaml"
	"github.com/jedib0t/go-pretty/table"
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

var listOpts = ListCSOpts{}
var showAll bool

func init() {
	getCmd.AddCommand(getCsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCsCmd.Flags().StringVarP(&listOpts.Tenant, "namespace-id", "n", "", "namespace id")
	getCsCmd.Flags().StringVarP(&listOpts.Group, "group", "g", "DEFAULT_GROUP", "group name")
	getCsCmd.Flags().IntVarP(&listOpts.PageNumber, "page-number", "P", 1, "page number")
	getCsCmd.Flags().IntVarP(&listOpts.PageSize, "page-size", "s", 10, "page size")
	getCsCmd.Flags().BoolVarP(&showAll, "all", "A", false, "show all configurations")
}

func GetCs(args []string) {
	naClient, err := NewNacosClient()
	if err != nil {
		log.Fatal(err)
	}
	configs, err := naClient.ListConfig(&listOpts)
	if err != nil {
		log.Fatal(err)
	}
	if len(args) > 0 {
		var items []*Config
		for _, ns := range configs.PageItems {
			if slices.Contains(args, ns.DataID) {
				items = append(items, ns)
			}
		}
		configs.PageItems = items
	}
	switch output {
	case "json":
		y, err := json.Marshal(configs)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(y))
	case "yaml":
		y, err := yaml.Marshal(configs)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(y))
	default:
		printCS(configs.PageItems)
	}

}

func printCS(css []*Config) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Data Id", "Group", "Application", "type"})
	for _, cs := range css {
		t.AppendRow(table.Row{cs.DataID, cs.Group, cs.AppName, cs.Type})
	}
	t.SortBy([]table.SortBy{{Name: "Group", Mode: table.Asc}, {Name: "Data Id", Mode: table.Asc}})

	s := table.StyleLight
	s.Options = table.OptionsNoBordersAndSeparators
	t.SetStyle(s)
	t.Render()
}
