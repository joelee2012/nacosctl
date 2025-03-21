/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/ghodss/yaml"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

var formats = []string{"table", "json", "yaml"}
var output string

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

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func PrintResources(v any, stdout io.Writer, format string) error {
	switch format {
	case "json":
		y, err := json.Marshal(v)
		if err != nil {
			return err
		}
		fmt.Fprint(stdout, string(y))
	case "yaml":
		y, err := yaml.Marshal(v)
		if err != nil {
			return err
		}
		fmt.Fprint(stdout, string(y))
	default:
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		switch v := v.(type) {
		case *NsList:
			t.AppendHeader(table.Row{"NAMESPACE", "ID", "Description", "Number"})
			for _, ns := range v.Items {
				t.AppendRow(table.Row{ns.ShowName, ns.Name, ns.Desc, ns.ConfigCount})
			}
			t.SortBy([]table.SortBy{{Name: "NAMESPACE", Mode: table.Asc}, {Name: "ID", Mode: table.Asc}})
		case *ConfigList:
			t.AppendHeader(table.Row{"DATAID", "GROUP", "TYPE", "CONTENT"})
			for _, cs := range v.PageItems {
				t.AppendRow(table.Row{cs.DataID, cs.Group, cs.AppName, cs.Type})
			}
			t.SortBy([]table.SortBy{{Name: "Group", Mode: table.Asc}, {Name: "Data Id", Mode: table.Asc}})
		}

		s := table.StyleLight
		s.Options = table.OptionsNoBordersAndSeparators
		t.SetStyle(s)
		t.Render()
	}
	return nil
}
