/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configdelCmd represents the configdel command
var configDelCmd = &cobra.Command{
	Use:     "del",
	Aliases: []string{"rm"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return delConfig(args[0])
	},
}

func init() {
	configCmd.AddCommand(configDelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configdelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configdelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func delConfig(name string) error {
	key := "servers." + name
	if !viper.IsSet(key) {
		return fmt.Errorf("%s not exists in %s", key, viper.ConfigFileUsed())
	}
	if viper.IsSet("context") && viper.GetString("context") == name {
		return fmt.Errorf("%s is current context, please use 'config set context' to set another context", name)
	}
	fmt.Printf("delete %s\n", key)

	var cliConfig CLIConfig
	viper.Unmarshal(&cliConfig)
	delete(cliConfig.Servers, name)
	data, err := json.MarshalIndent(cliConfig, "", "  ")
	if err != nil {
		return err
	}
	err = viper.ReadConfig(bytes.NewReader(data))
	if err != nil {
		return err
	}
	return viper.WriteConfig()
}
