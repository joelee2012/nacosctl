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

	"github.com/joelee2012/nacosctl/pkg/nacos"
	"github.com/spf13/cobra"
)

// getCsCmd represents the getCs command
var getCsCmd = &cobra.Command{
	Use:     "cs [name]",
	Aliases: []string{"configuration"},
	Short:   "Display one or many configurations",
	Run: func(cmd *cobra.Command, args []string) {
		GetCs(args)
	},
}

func init() {
	getCmd.AddCommand(getCsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCsCmd.Flags().StringVarP(&cmdOpts.NamespaceID, "namespace", "n", "", "namespace id")
	getCsCmd.Flags().StringVarP(&cmdOpts.Group, "group", "g", "DEFAULT_GROUP", "group name")
	getCsCmd.Flags().BoolVarP(&cmdOpts.ShowAll, "all", "A", false, "show all configurations")

}

func GetCs(args []string) {
	client := NewNacosClient()
	allCs := new(nacos.ConfigurationList)
	var err error
	if cmdOpts.ShowAll {
		allCs, err = client.ListAllConfig()
		cobra.CheckErr(err)
	} else {
		if len(args) > 0 {

			for _, c := range args {
				cs, err := client.GetConfig(&nacos.GetCfgOpts{NamespaceID: cmdOpts.NamespaceID, Group: cmdOpts.Group, DataID: c})
				cobra.CheckErr(err)
				allCs.Items = append(allCs.Items, *cs)
			}
		} else {
			allCs, err = client.ListConfigInNs(cmdOpts.NamespaceID, cmdOpts.Group)
			cobra.CheckErr(err)
		}
	}
	list := NewList(client.APIVersion, allCs.Items, NewConfiguration)
	// toJson(list, os.Stdout)
	cobra.CheckErr(WriteFormat(list, cmdOpts.Output, os.Stdout))
}
