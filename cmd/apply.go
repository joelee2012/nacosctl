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
	if err := ns.LoadFromYaml(name); err == nil {
		cobra.CheckErr(client.CreateOrUpdateNamespace(&nacos.CreateNSOpts{ID: ns.Name, Desc: ns.Desc, Name: ns.ShowName}))
		fmt.Printf("namespace/%s created\n", ns.ShowName)
		return
	}
	nsNames := ListNamespace(client)
	c := &nacos.Config{}
	cobra.CheckErr(c.LoadFromYaml(name))
	if !slices.Contains(nsNames, c.Tenant) {
		cobra.CheckErr(fmt.Errorf("namespace/%s not found", c.Tenant))
	}
	cobra.CheckErr(client.CreateConfig(&nacos.CreateCSOpts{
		DataID:  c.DataID,
		Group:   c.Group,
		Tenant:  c.Tenant,
		Content: c.Content,
		Type:    c.Type,
		Desc:    c.Desc,
		AppName: c.AppName,
		Tags:    c.Tags,
	}))
	fmt.Printf("configuration/%s created\n", c.DataID)
}

func ListNamespace(client *nacos.Client) []string {
	nsList, err := client.ListNamespace()
	cobra.CheckErr(err)
	nsNames := []string{}
	for _, ns := range nsList.Items {
		nsNames = append(nsNames, ns.Name)
	}
	return nsNames
}
func CreateResourceFromDir(naClient *nacos.Client, dir string) {
	nss := new(nacos.NsList)
	cs := new(nacos.ConfigList)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		cobra.CheckErr(err)
		if !info.IsDir() {
			ns := &nacos.Namespace{}
			if err := ns.LoadFromYaml(path); err == nil {
				nss.Items = append(nss.Items, ns)
			} else {
				c := &nacos.Config{}
				cobra.CheckErr(c.LoadFromYaml(path))
				cs.Items = append(cs.Items, c)
			}
		}
		return nil
	})
	cobra.CheckErr(err)
	nsNames := ListNamespace(naClient)
	for _, ns := range nss.Items {
		cobra.CheckErr(naClient.CreateOrUpdateNamespace(&nacos.CreateNSOpts{ID: ns.Name, Desc: ns.Desc, Name: ns.ShowName}))
		fmt.Printf("namespace/%s created\n", ns.ShowName)
		if !slices.Contains(nsNames, ns.ShowName) {
			nsNames = append(nsNames, ns.ShowName)
		}
	}
	for _, c := range cs.Items {
		if !slices.Contains(nsNames, c.Tenant) {
			cobra.CheckErr(fmt.Errorf("namespace/%s not found", c.Tenant))
		}
		cobra.CheckErr(naClient.CreateConfig(&nacos.CreateCSOpts{
			DataID:  c.DataID,
			Group:   c.Group,
			Tenant:  c.Tenant,
			Content: c.Content,
			Type:    c.Type,
			Desc:    c.Desc,
			AppName: c.AppName,
			Tags:    c.Tags,
		}))
		fmt.Printf("configuration/%s created\n", c.DataID)
	}
}
