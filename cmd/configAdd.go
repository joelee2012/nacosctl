/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configAddCmd represents the configAdd command
var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new nacos server",
	Run: func(cmd *cobra.Command, args []string) {
		AddConfig(args[0])
	},
	Args: cobra.ExactArgs(1),
}

var na = &Nacos{}

func init() {
	configCmd.AddCommand(configAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	configAddCmd.Flags().StringVar(&na.URL, "url", "", "the nacos url")
	configAddCmd.MarkFlagRequired("url")
	configAddCmd.Flags().StringVarP(&na.User, "user", "u", "", "nacos user")
	configAddCmd.MarkFlagRequired("user")
	configAddCmd.Flags().StringVarP(&na.Password, "password", "p", "", "nacos password")
	configAddCmd.MarkFlagRequired("password")
}

func AddConfig(name string) {
	if viper.IsSet("servers." + name) {
		fmt.Printf("server.%s already exists in %s\n", name, viper.ConfigFileUsed())
		return
	} else {
		server := Server{
			URL:      na.URL,
			User:     na.User,
			Password: na.Password,
		}
		viper.Set("servers."+name, server)
		if !viper.IsSet("context") {
			viper.Set("context", name)
		}
	}

	viper.WriteConfig()
}
