package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/joelee2012/nacosctl/pkg/nacos"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var applyCmd = &cobra.Command{
	Use:   "apply [flags] [command]",
	Short: "Apply configuration file to nacos",
	Run: func(cmd *cobra.Command, args []string) {
		if cmdOpts.OutDir != "" {
			client := NewNacosClient()
			fi, err := os.Stat(cmdOpts.OutDir)
			cobra.CheckErr(err)
			switch mode := fi.Mode(); {
			case mode.IsRegular():
				CreateResourceFromFile(client, cmdOpts.OutDir)
			case mode.IsDir():
				CreateResourceFromDir(client, cmdOpts.OutDir)
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
	applyCmd.Flags().StringVarP(&cmdOpts.OutDir, "filename", "f", "", "The files or dir that contain the configurations")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func CreateResourceFromFile(client *nacos.Client, name string) {
	ns := new(Namespace)
	if err := readYamlFile(ns, name); err == nil {
		cobra.CheckErr(client.CreateOrUpdateNamespace(&nacos.CreateNSOpts{ID: ns.Metadata.ID, Description: ns.Metadata.Description, Name: ns.Metadata.Name}))
		fmt.Printf("namespace/%s created\n", ns.Metadata.Name)
		return
	}
	nsNames := ListNamespace(client)
	c := new(Configuration)
	cobra.CheckErr(readYamlFile(c, name))
	if !slices.Contains(nsNames, c.Metadata.Namespace) {
		cobra.CheckErr(fmt.Errorf("namespace/%s not found", c.Metadata.Namespace))
	}
	cobra.CheckErr(client.CreateConfig(&nacos.CreateCSOpts{
		DataID:      c.Metadata.DataID,
		Group:       c.Metadata.Group,
		NamespaceID: c.Metadata.Namespace,
		Content:     c.Spec.Content,
		Type:        c.Spec.Type,
		Description: c.Spec.Description,
		Application: c.Spec.Application,
		Tags:        c.Spec.Tags,
	}))
	fmt.Printf("configuration/%s created\n", c.Metadata.DataID)
}

func ListNamespace(client *nacos.Client) []string {
	nsList, err := client.ListNamespace()
	cobra.CheckErr(err)
	nsNames := []string{}
	for _, ns := range nsList.Items {
		nsNames = append(nsNames, ns.ID)
	}
	return nsNames
}
func CreateResourceFromDir(naClient *nacos.Client, dir string) {
	nss := new(NamespaceList)
	cs := new(ConfigurationList)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		cobra.CheckErr(err)
		if !info.IsDir() {
			var ns Namespace
			if err := readYamlFile(ns, path); err == nil {
				nss.Items = append(nss.Items, ns)
			} else {
				var c Configuration
				cobra.CheckErr(readYamlFile(c, path))
				cs.Items = append(cs.Items, c)
			}
		}
		return nil
	})
	cobra.CheckErr(err)
	nsNames := ListNamespace(naClient)
	for _, ns := range nss.Items {
		cobra.CheckErr(naClient.CreateOrUpdateNamespace(&nacos.CreateNSOpts{ID: ns.Metadata.ID, Description: ns.Metadata.Description, Name: ns.Metadata.Name}))
		fmt.Printf("namespace/%s created\n", ns.Metadata.Name)
		if !slices.Contains(nsNames, ns.Metadata.Name) {
			nsNames = append(nsNames, ns.Metadata.Name)
		}
	}
	for _, c := range cs.Items {
		if !slices.Contains(nsNames, c.Metadata.Namespace) {
			cobra.CheckErr(fmt.Errorf("namespace/%s not found", c.Metadata.Namespace))
		}
		cobra.CheckErr(naClient.CreateConfig(&nacos.CreateCSOpts{
			DataID:      c.Metadata.DataID,
			Group:       c.Metadata.Group,
			NamespaceID: c.Metadata.Namespace,
			Content:     c.Spec.Content,
			Type:        c.Spec.Type,
			Description: c.Spec.Description,
			Application: c.Spec.Application,
			Tags:        c.Spec.Tags,
		}))
		fmt.Printf("configuration/%s created\n", c.Metadata.DataID)
	}
}
