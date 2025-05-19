/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage nacos instance config",
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type CLIConfig struct {
	Context string             `json:"context"`
	Servers map[string]*Server `json:"servers"`
}

type Server struct {
	Password string `json:"password"`
	URL      string `json:"url"`
	User     string `json:"user"`
}

func (c *CLIConfig) ReadFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(c)
}

func (c *CLIConfig) WriteFile(name string) error {
	return writeFile(c, name)
}

func (c *CLIConfig) GetServer(name string) *Server {
	return c.Servers[name]
}

func (c *CLIConfig) AddServer(name string, server *Server) {
	c.Servers[name] = server
}

func (c *CLIConfig) DeleteServer(name string) {
	delete(c.Servers, name)
	if c.Context == name {
		c.Context = ""
	}
}

func (c *CLIConfig) SetContext(name string) error {
	for k := range c.Servers {
		if k == name {
			c.Context = name
			return nil
		}
	}
	return fmt.Errorf("server %s not found", name)
}

func (c *CLIConfig) GetContext() string {
	return c.Context
}

func (c *CLIConfig) GetCurrentServer() *Server {
	return c.Servers[c.Context]
}

func (c *CLIConfig) ToYaml() ([]byte, error) {
	return yaml.Marshal(c)
}
