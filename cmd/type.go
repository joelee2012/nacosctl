package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/jedib0t/go-pretty/table"
	"github.com/joelee2012/nacosctl/pkg/nacos"
)

type FormatWriter interface {
	TableWriter
	JsonWriter
	YamlWriter
}

type TableWriter interface {
	ToTable(w io.Writer)
}

type JsonWriter interface {
	ToJson(w io.Writer) error
}
type YamlWriter interface {
	ToYaml(w io.Writer) error
}

type FileWriter interface {
	ToFile(w io.Writer) error
}
type DirWriter interface {
	WriteToDir(name string) error
}

// YamlFileLoader interface for loading from YAML
type YamlFileLoader interface {
	FromYaml(name string) error
}

func toJson(v any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func toYaml(v any, w io.Writer) error {
	enc := yaml.NewEncoder(w)
	return enc.Encode(v)
}

func readYamlFile(v any, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := yaml.NewDecoder(f, yaml.DisallowUnknownField())
	return dec.Decode(v)
}

func writeYamlFile(v any, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return toYaml(v, f)
}

func toTable(w io.Writer, fn func(t table.Writer)) {
	tb := table.NewWriter()
	tb.SetOutputMirror(w)
	fn(tb)
	s := table.StyleLight
	s.Options = table.OptionsNoBordersAndSeparators
	tb.SetStyle(s)
	tb.Render()
}

func WriteAsFormat(format string, writable FormatWriter, w io.Writer) {
	switch format {
	case "json":
		writable.ToJson(w)
	case "yaml":
		writable.ToYaml(w)
	case "table":
		writable.ToTable(w)
	default:
		writable.ToTable(w)
	}
}

func WriteToDir(name string, writable DirWriter) error {
	return writable.WriteToDir(name)
}
func LoadFromYaml(name string, loader YamlFileLoader) error {
	return loader.FromYaml(name)
}

type Configuration struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Application string `json:"application"`
		Group       string `json:"group"`
		DataID      string `json:"name"`
		Namespace   string `json:"namespace"`
		Type        string `json:"type"`
		Tags        string `json:"tags"`
		Description string `json:"description"`
	} `json:"metadata"`
	Data   string `json:"data"`
	Status struct {
		Md5              string `json:"md5,omitempty"`
		EncryptedDataKey string `json:"encryptedDataKey,omitempty"`
		CreateTime       int64  `json:"createTime,omitempty"`
		ModifyTime       int64  `json:"modifyTime,omitempty"`
	} `json:"status"`
}

func (c *Configuration) ToJson(w io.Writer) error {
	return toJson(c, w)
}

func (c *Configuration) ToYaml(w io.Writer) error {
	return toYaml(c, w)
}

func (c *Configuration) ToFile(name string) error {
	return writeYamlFile(c, name)
}

// FromYaml load from YAML file
func (c *Configuration) FromYaml(name string) error {
	return readYamlFile(c, name)
}

func (c *Configuration) FromNacosConfig(apiVersion string, nc *nacos.Config) {
	c.APIVersion = apiVersion
	c.Kind = "Configuration"
	c.Metadata.DataID = nc.DataID
	c.Metadata.Namespace = nc.NamespaceID
	c.Metadata.Group = nc.Group
	c.Metadata.Application = nc.Application
	c.Metadata.Type = nc.Type
	c.Metadata.Tags = nc.Tags
	c.Metadata.Description = nc.Description
	c.Data = nc.Content
	c.Status.Md5 = nc.Md5
	c.Status.CreateTime = nc.CreateTime
	c.Status.ModifyTime = nc.ModifyTime
	c.Status.EncryptedDataKey = nc.EncryptedDataKey
}

type ConfigurationList struct {
	APIVersion string           `json:"apiVersion"`
	Items      []*Configuration `json:"items"`
	Kind       string           `json:"kind"`
}

func (list *ConfigurationList) ToTable(w io.Writer) {
	toTable(w, func(t table.Writer) {
		t.AppendHeader(table.Row{"NAMESPACEID", "DATAID", "GROUP", "APPLICATION", "TYPE"})
		for _, item := range list.Items {
			t.AppendRow(table.Row{item.Metadata.Namespace, item.Metadata.DataID, item.Metadata.Group, item.Metadata.Application, item.Metadata.Type})
		}
		t.SortBy([]table.SortBy{{Name: "NAMESPACEID", Mode: table.Asc}, {Name: "DATAID", Mode: table.Asc}})
	})
}

func (list *ConfigurationList) ToJson(w io.Writer) error {
	return toJson(list, w)
}

func (list *ConfigurationList) ToYaml(w io.Writer) error {
	return toYaml(list, w)
}

func (list *ConfigurationList) FromNacosConfigList(apiVersion string, cs *nacos.ConfigList) {
	list.Kind = "List"
	list.APIVersion = apiVersion
	for _, item := range cs.Items {
		var c Configuration
		c.FromNacosConfig("v1", item)
		list.Items = append(list.Items, &c)
	}
}

func (list *ConfigurationList) WriteToDir(name string) error {
	var dir string
	for _, c := range list.Items {
		if c.Metadata.Namespace == "" {
			dir = filepath.Join(name, "public", c.Metadata.Group)
		} else {
			dir = filepath.Join(name, c.Metadata.Namespace, c.Metadata.Group)
		}
		if err := os.MkdirAll(dir, 0750); err != nil {
			return err
		}
		if err := c.ToFile(filepath.Join(dir, c.Metadata.DataID)); err != nil {
			return err
		}
	}
	return nil
}

type NamespaceList struct {
	APIVersion string       `json:"apiVersion"`
	Items      []*Namespace `json:"items"`
	Kind       string       `json:"kind"`
}

func (list *NamespaceList) ToJson(w io.Writer) error {
	return toJson(list, w)
}

func (list *NamespaceList) ToYaml(w io.Writer) error {
	return toYaml(list, w)
}

func (list *NamespaceList) FromNacosNsList(apiVersion string, ns *nacos.NsList) {
	list.Kind = "List"
	list.APIVersion = apiVersion
	for _, item := range ns.Items {
		var n Namespace
		n.FromNacosNamespace("v1", item)
		list.Items = append(list.Items, &n)
	}
}

func (list *NamespaceList) ToTable(w io.Writer) {
	toTable(w, func(t table.Writer) {
		t.AppendHeader(table.Row{"NAME", "ID", "DESCRIPTION", "COUNT"})
		for _, ns := range list.Items {
			t.AppendRow(table.Row{ns.Metadata.Name, ns.Metadata.ID, ns.Metadata.Description, ns.Status.ConfigCount})
		}
		t.SortBy([]table.SortBy{{Name: "NAME", Mode: table.Asc}, {Name: "ID", Mode: table.Asc}})
	})
}

func (n *NamespaceList) WriteToDir(name string) error {
	for _, e := range n.Items {
		if e.Metadata.ID == "" {
			continue
		}
		if err := e.ToFile(filepath.Join(name, fmt.Sprintf("%s.yaml", e.Metadata.ID))); err != nil {
			return err
		}
	}
	return nil
}

type Namespace struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name        string `json:"name"`
		ID          string `json:"id"`
		Description string `json:"description"`
	} `json:"metadata"`
	Status struct {
		Quota       int `json:"quota"`
		ConfigCount int `json:"configCount"`
		Type        int `json:"type"`
	} `json:"status"`
}

func (n *Namespace) ToJson(w io.Writer) error {
	return toJson(n, w)
}

func (n *Namespace) ToYaml(w io.Writer) error {
	return toYaml(n, w)
}
func (n *Namespace) ToFile(name string) error {
	return writeYamlFile(n, name)
}

func (n *Namespace) FromYaml(name string) error {
	return readYamlFile(n, name)
}

func (n *Namespace) FromNacosNamespace(apiVersion string, e *nacos.Namespace) {
	n.APIVersion = apiVersion
	n.Kind = "Namespace"
	n.Metadata.Name = e.Name
	n.Metadata.ID = e.ID
	n.Metadata.Description = e.Description
	n.Status.ConfigCount = e.ConfigCount
	n.Status.Quota = e.Quota
	n.Status.Type = e.Type
}
