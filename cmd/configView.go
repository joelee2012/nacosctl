/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configViewCmd represents the configView command
var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View nacos config",

	Run: func(cmd *cobra.Command, args []string) {
		c := viper.AllSettings()
		bs, err := yaml.Marshal(c)
		if err != nil {
			log.Fatalf("unable to marshal config to YAML: %v", err)
		}
		fmt.Println(string(bs))
	},
}

func init() {
	configCmd.AddCommand(configViewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configViewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configViewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
