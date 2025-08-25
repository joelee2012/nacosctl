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
	ns := &nacos.Namespace{}
	if err := ns.FromYaml(name); err == nil {
		cobra.CheckErr(client.CreateOrUpdateNamespace(&nacos.CreateNSOpts{ID: ns.ID, Description: ns.Description, Name: ns.Name}))
		fmt.Printf("namespace/%s created\n", ns.Name)
		return
	}
	nsNames := ListNamespace(client)
	c := &nacos.Configuration{}
	cobra.CheckErr(c.FromYaml(name))
	if !slices.Contains(nsNames, c.NamespaceID) {
		cobra.CheckErr(fmt.Errorf("namespace/%s not found", c.NamespaceID))
	}
	cobra.CheckErr(client.CreateConfig(&nacos.CreateCSOpts{
		DataID:      c.DataID,
		Group:       c.Group,
		NamespaceID: c.NamespaceID,
		Content:     c.Content,
		Type:        c.Type,
		Description: c.Description,
		Application: c.Application,
		Tags:        c.Tags,
	}))
	fmt.Printf("configuration/%s created\n", c.DataID)
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
	nss := new(nacos.NamespaceList)
	cs := new(nacos.ConfigurationList)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		cobra.CheckErr(err)
		if !info.IsDir() {
			ns := &nacos.Namespace{}
			if err := ns.FromYaml(path); err == nil {
				nss.Items = append(nss.Items, ns)
			} else {
				c := &nacos.Configuration{}
				cobra.CheckErr(c.FromYaml(path))
				cs.Items = append(cs.Items, c)
			}
		}
		return nil
	})
	cobra.CheckErr(err)
	nsNames := ListNamespace(naClient)
	for _, ns := range nss.Items {
		cobra.CheckErr(naClient.CreateOrUpdateNamespace(&nacos.CreateNSOpts{ID: ns.ID, Description: ns.Description, Name: ns.Name}))
		fmt.Printf("namespace/%s created\n", ns.Name)
		if !slices.Contains(nsNames, ns.Name) {
			nsNames = append(nsNames, ns.Name)
		}
	}
	for _, c := range cs.Items {
		if !slices.Contains(nsNames, c.NamespaceID) {
			cobra.CheckErr(fmt.Errorf("namespace/%s not found", c.NamespaceID))
		}
		cobra.CheckErr(naClient.CreateConfig(&nacos.CreateCSOpts{
			DataID:      c.DataID,
			Group:       c.Group,
			NamespaceID: c.NamespaceID,
			Content:     c.Content,
			Type:        c.Type,
			Description: c.Description,
			Application: c.Application,
			Tags:        c.Tags,
		}))
		fmt.Printf("configuration/%s created\n", c.DataID)
	}
}
