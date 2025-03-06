/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// var client *nacos.nacos

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nactl",
	Short: "Command line tools for Nacos",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// type NacosOpts struct {
// 	URL      string
// 	User     string
// 	Password string
// 	Output   string
// }

// var naClient = &Nacos{}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nacos.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cmd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".nacos")
	}
	viper.SetEnvPrefix("nacos")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			cobra.CheckErr(err)
		}
	}
}

func NewNacosClient() (*Nacos, error) {
	var n Nacos
	if viper.IsSet("context") {
		ctx := viper.GetString("context")
		key := "servers." + ctx
		if v := viper.Sub(key); v != nil {
			err := v.Unmarshal(&n)
			if err != nil {
				return nil, err
			}
			if n.URL == "" {
				return nil, fmt.Errorf("%s no url set", key)
			}
			if n.User == "" {
				return nil, fmt.Errorf("%s no user set", key)
			}
			if n.Password == "" {
				return nil, fmt.Errorf("%s no password set", key)
			}
			return &n, nil
		}
	}
	return nil, fmt.Errorf("no context set")
}
