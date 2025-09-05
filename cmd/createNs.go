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

// createNsCmd represents the createNs command
var createNsCmd = &cobra.Command{
	Use:   "ns name",
	Short: "Create one namespace",
	Run: func(cmd *cobra.Command, args []string) {
		client := NewNacosClient()
		nsOpts.Name = args[0]
		cobra.CheckErr(client.CreateNamespace(&nsOpts))
		fmt.Printf("namespace/%#v created\n", nsOpts)
	},
	Args: cobra.ExactArgs(1),
}

var nsOpts nacos.CreateNsOpts

func init() {
	createCmd.AddCommand(createNsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createNsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createNsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// createNsCmd.Flags().StringVarP(&nsOpts.Name, "name", "n", "", "name of namespace")
	createNsCmd.Flags().StringVarP(&nsOpts.ID, "id", "i", "", "id of namespace")
	createNsCmd.MarkFlagRequired("id")
	createNsCmd.Flags().StringVarP(&nsOpts.Description, "desc", "d", "", "description of namespace")
	createNsCmd.MarkFlagRequired("desc")
}
