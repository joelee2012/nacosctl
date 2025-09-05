/*
Copyright © 2025 Joe Lee <lj_2005@163.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/joelee2012/nacosctl/pkg/nacos"
	"github.com/spf13/cobra"
)

// createCsCmd represents the createCs command
var createCsCmd = &cobra.Command{
	Use:   "cs [flags] name",
	Short: "Create one configuration",

	Run: func(cmd *cobra.Command, args []string) {
		client := NewNacosClient()
		createOpts.DataID = args[0]
		cobra.CheckErr(client.CreateConfig(&createOpts))
		fmt.Printf("configuration/%s/%s/%s created\n", createOpts.NamespaceID, createOpts.Group, createOpts.DataID)

	},
	Args: cobra.ExactArgs(1),
}

var createOpts nacos.CreateCfgOpts

func init() {
	createCmd.AddCommand(createCsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCsCmd.Flags().StringVarP(&createOpts.NamespaceID, "namespace", "n", "", "namespace id")
	createCsCmd.Flags().StringVarP(&createOpts.Group, "group", "g", "DEFAULT_GROUP", "group of configuration")
	createCsCmd.Flags().StringVarP(&createOpts.Content, "content", "c", "", "content of configuration")
	createCsCmd.MarkFlagRequired("content")
	createCsCmd.Flags().StringVarP(&createOpts.Type, "type", "t", "text", "configuration type")
	createCsCmd.Flags().StringVarP(&createOpts.Description, "description", "d", "", "description of configuration")
	createCsCmd.Flags().StringVarP(&createOpts.Tags, "tags", "T", "", "tags of configuration")
	createCsCmd.Flags().StringVarP(&createOpts.Application, "application", "a", "", "application of configuration")

}
