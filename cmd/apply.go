package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var applyCmd = &cobra.Command{
	Use:   "apply  [flags] [command]",
	Short: "Apply configuration file  to nacos",
	Run: func(cmd *cobra.Command, args []string) {
		if cmdOpts.OutDir != "" {
			naClient, err := NewNacosClient()
			cobra.CheckErr(err)
			if IsFile(cmdOpts.OutDir) {
				ns := &Namespace{}
				cobra.CheckErr(ns.FromYaml(cmdOpts.OutDir))
				cobra.CheckErr(naClient.CreateNamespace(&CreateNSOpts{ID: ns.Name, Desc: ns.Desc, Name: ns.ShowName}))
				fmt.Printf("namespace/%s created\n", ns.ShowName)
			}
			if IsDir(cmdOpts.OutDir) {
				ns := &Namespace{}
				err = filepath.Walk(cmdOpts.OutDir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
						return err
					}
					if !info.IsDir() {
						if err = ns.FromYaml(path); err != nil {
							return err
						}
						if err = naClient.CreateNamespace(&CreateNSOpts{ID: ns.Name, Desc: ns.Desc, Name: ns.ShowName}); err != nil {
							return err
						}
						fmt.Printf("namespace/%s created\n", ns.ShowName)
					}
					return nil
				})
				cobra.CheckErr(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")
	applyCmd.Flags().StringVarP(&cmdOpts.OutDir, "filename", "f", "", "The files that contain the configurations")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
