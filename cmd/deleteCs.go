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

// deleteCsCmd represents the deleteCs command
var deleteCsCmd = &cobra.Command{
	Use:     "cs",
	Aliases: []string{"configuration"},
	Short:   "Delete one or many configurations",
	Run: func(cmd *cobra.Command, args []string) {
		client := NewNacosClient()
		for _, dataId := range args {
			err := client.DeleteConfig(&nacos.DeleteCfgOpts{
				DataID:      dataId,
				Group:       cmdOpts.Group,
				NamespaceID: cmdOpts.NamespaceID,
			})
			cobra.CheckErr(err)
			fmt.Printf("configuration/%s deleted\n", dataId)

		}
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	deleteCmd.AddCommand(deleteCsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deleteCsCmd.Flags().StringVarP(&cmdOpts.NamespaceID, "namespace", "n", "", "namespace id")
	deleteCsCmd.Flags().StringVarP(&cmdOpts.Group, "group", "g", "DEFAULT_GROUP", "name of group")
}
