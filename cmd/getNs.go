/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/ghodss/yaml"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

// getNsCmd represents the getNs command
var getNsCmd = &cobra.Command{
	Use:     "namespace [name]",
	Aliases: []string{"ns"},
	Short:   "Display one or many namespaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		return GetNamespace(args)
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

func GetNamespace(args []string) error {
	naClient, err := NewNacosClient()
	if err != nil {
		return err
	}
	nss, err := naClient.ListNamespace()
	if err != nil {
		return err
	}
	if len(args) > 0 {
		var items []*Namespace
		for _, ns := range nss.Items {
			if slices.Contains(args, ns.ShowName) {
				items = append(items, ns)
			}
		}
		nss.Items = items
	}
	switch output {
	case "json":
		y, err := json.Marshal(nss)
		if err != nil {
			return err
		}
		fmt.Println(string(y))
	case "yaml":
		y, err := yaml.Marshal(nss)
		if err != nil {
			return err
		}
		fmt.Println(string(y))
	default:
		printTable(nss.Items)
	}

	// marshalFunc := json.Marshal

	// if output == "yaml" {
	// 	marshalFunc = yaml.Marshal
	// }
	// y, err := marshalFunc(nss)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(string(y))
	return nil
}

func printTable(nss []*Namespace) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"NAMESPACE", "ID", "Description", "Number"})
	for _, ns := range nss {
		t.AppendRow(table.Row{ns.ShowName, ns.Name, ns.Desc, ns.ConfigCount})
	}
	t.SortBy([]table.SortBy{{Name: "NAMESPACE", Mode: table.Asc}, {Name: "ID", Mode: table.Asc}})

	s := table.StyleLight
	s.Options = table.OptionsNoBordersAndSeparators
	t.SetStyle(s)
	t.Render()
}
