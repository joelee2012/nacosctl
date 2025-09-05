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

// getUserCmd represents the getCs command
var getUserCmd = &cobra.Command{
	Use:     "user [name]",
	Aliases: []string{"u"},
	Short:   "Display one or many user",
	Run: func(cmd *cobra.Command, args []string) {
		getUser(args)
	},
}

func init() {
	getCmd.AddCommand(getUserCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getUserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

}

func getUser(args []string) {
	client := NewNacosClient()
	users, err := client.ListUser()
	cobra.CheckErr(err)
	if len(args) > 0 {
		var us []*nacos.User
		for _, u := range users.Items {
			if slices.Contains(args, u.Name) {
				us = append(us, u)
			}
		}
		users.Items = us
	}
	list := NewList(client.APIVersion, users.Items, NewUser)
	cobra.CheckErr(WriteFormat(list, cmdOpts.Output, os.Stdout))
}
