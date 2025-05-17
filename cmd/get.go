/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"slices"

	"github.com/ghodss/yaml"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

var formats = []string{"table", "json", "yaml"}
var output string
var outputDir string

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many resources",
	Long:  `Prints a table of information about the specified resources`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")

		if !slices.Contains(formats, output) {
			return fmt.Errorf("invalid format type [%s] for [-o, --output] flag", output)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	getCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "A help for foo")
	getCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "O", "", "output directory")
	getCmd.MarkFlagsMutuallyExclusive("output", "output-dir")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func printJson(v any, stdout io.Writer) error {
	y, err := json.Marshal(v)
	if err != nil {
		return err
	}
	fmt.Fprint(stdout, string(y))
	return nil
}

func printYaml(v any, stdout io.Writer) error {
	y, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	fmt.Fprint(stdout, string(y))
	return nil
}

func saveToFile(v any, folder string) error {
	switch v := v.(type) {
	case *NsList:
		for _, ns := range v.Items {
			data, err := yaml.Marshal(ns)
			if err != nil {
				return err
			}
			if err := os.WriteFile(path.Join(folder, ns.ShowName), data, 0666); err != nil {
				return err
			}
		}
	case *ConfigList:
		for _, cs := range v.PageItems {
			data, err := yaml.Marshal(cs)
			if err != nil {
				return err
			}
			if err := os.WriteFile(path.Join(folder, cs.DataID), data, 0666); err != nil {
				return err
			}
		}
	}
	return nil
}

func printItem(v any, stdout io.Writer, format string) error {
	funcMap := map[string]func(any, io.Writer) error{
		"json": printJson,
		"yaml": printYaml,
	}
	switch v := v.(type) {
	case *NsList:
		if len(v.Items) == 1 {
			return funcMap[format](v.Items[0], stdout)
		} else {
			return funcMap[format](v, stdout)
		}
	case *ConfigList:
		if len(v.PageItems) == 1 {
			return funcMap[format](v.PageItems[0], stdout)
		} else {
			return funcMap[format](v, stdout)
		}
	}
	return nil
}
func PrintResources(v any, stdout io.Writer, format string) error {
	if outputDir != "" {
		return saveToFile(v, outputDir)
	}
	switch format {
	case "json", "yaml":
		printItem(v, stdout, format)
	default:
		t := table.NewWriter()
		t.SetOutputMirror(stdout)
		switch v := v.(type) {
		case *NsList:
			t.AppendHeader(table.Row{"NAMESPACE", "ID", "DESCRIPTION", "COUNT"})
			for _, ns := range v.Items {
				t.AppendRow(table.Row{ns.ShowName, ns.Name, ns.Desc, ns.ConfigCount})
			}
			t.SortBy([]table.SortBy{{Name: "NAMESPACE", Mode: table.Asc}, {Name: "ID", Mode: table.Asc}})
		case *ConfigList:
			t.AppendHeader(table.Row{"NAMESPACE", "DATAID", "GROUP", "APPLICATION", "TYPE"})
			for _, cs := range v.PageItems {
				t.AppendRow(table.Row{cs.Tenant, cs.DataID, cs.Group, cs.AppName, cs.Type})
			}
			t.SortBy([]table.SortBy{{Name: "NAMESPACE", Mode: table.Asc}, {Name: "DATAID", Mode: table.Asc}})
		}

		s := table.StyleLight
		s.Options = table.OptionsNoBordersAndSeparators
		t.SetStyle(s)
		t.Render()
	}
	return nil
}
