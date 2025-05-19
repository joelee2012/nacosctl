/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

var formats = []string{"table", "json", "yaml"}

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
	getCmd.PersistentFlags().StringVarP(&cmdOpts.Output, "output", "o", "table", "Output format. One of: (json, yaml, table]")
	getCmd.PersistentFlags().StringVarP(&cmdOpts.OutDir, "output-dir", "d", "", "Output directory")
	getCmd.MarkFlagsMutuallyExclusive("output", "output-dir")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
