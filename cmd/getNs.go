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
	"os"
	"slices"

	"github.com/joelee2012/nacosctl/pkg/nacos"
	"github.com/spf13/cobra"
)

// getNsCmd represents the getNs command
var getNsCmd = &cobra.Command{
	Use:     "ns [name]",
	Aliases: []string{"namespace"},
	Short:   "Display one or many namespaces",
	Run: func(cmd *cobra.Command, args []string) {
		GetNamespace(args)
	},
}

func init() {
	getCmd.AddCommand(getNsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getNsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getNsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GetNamespace(args []string) {
	client := NewNacosClient()
	nss, err := client.ListNamespace()
	cobra.CheckErr(err)
	if len(args) > 0 {
		var items []nacos.Namespace
		for _, ns := range nss.Items {
			if slices.Contains(args, ns.ID) {
				items = append(items, ns)
			}
		}
		nss.Items = items
	}
	list := NewList(client.APIVersion, nss.Items, NewNamespace)
	cobra.CheckErr(WriteFormat(list, cmdOpts.Output, os.Stdout))
}
