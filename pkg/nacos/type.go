package nacos

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jedib0t/go-pretty/table"
)

type FormatWriter interface {
	TableWriter
	JsonWriter
	YamlWriter
}

type TableWriter interface {
	WriteTable(w io.Writer)
}

type JsonWriter interface {
	WriteJson(w io.Writer) error
}
type YamlWriter interface {
	WriteYaml(w io.Writer) error
}

type FileWriter interface {
	WriteFile(w io.Writer) error
}
type DirWriter interface {
	WriteToDir(name string) error
}

// YamlFileLoader interface for loading from YAML
type YamlFileLoader interface {
	LoadFromYaml(name string) error
}

type ConfigList struct {
	TotalCount     int       `json:"totalCount,omitempty"`
	PageNumber     int       `json:"pageNumber,omitempty"`
	PagesAvailable int       `json:"pagesAvailable,omitempty"`
	Items          []*Config `json:"pageItems"`
}

func (c *ConfigList) WriteTable(w io.Writer) {
	writeTable(w, func(t table.Writer) {
		t.AppendHeader(table.Row{"NAMESPACE", "DATAID", "GROUP", "APPLICATION", "TYPE"})
		for _, item := range c.Items {
			t.AppendRow(table.Row{item.Tenant, item.DataID, item.Group, item.AppName, item.Type})
		}
		t.SortBy([]table.SortBy{{Name: "NAMESPACE", Mode: table.Asc}, {Name: "DATAID", Mode: table.Asc}})
	})
}

func (c *ConfigList) WriteJson(w io.Writer) error {
	return writeJson(c, w)
}

func (c *ConfigList) WriteYaml(w io.Writer) error {
	return writeYaml(c, w)
}

func (cs *ConfigList) WriteToDir(name string) error {
	var dir string
	for _, c := range cs.Items {
		if c.Tenant == "" {
			dir = filepath.Join(name, "public", c.Group)
		} else {
			dir = filepath.Join(name, c.Tenant, c.Group)
		}
		if err := os.MkdirAll(dir, 0750); err != nil {
			return err
		}
		if err := c.WriteFile(filepath.Join(dir, c.DataID)); err != nil {
			return err
		}
	}
	return nil
}

type Config struct {
	ID               string `json:"id"`
	DataID           string `json:"dataId"`
	Group            string `json:"group"`
	Content          string `json:"content"`
	Tenant           string `json:"tenant"`
	Type             string `json:"type"`
	Md5              string `json:"md5,omitempty"`
	EncryptedDataKey string `json:"encryptedDataKey,omitempty"`
	AppName          string `json:"appName,omitempty"`
	CreateTime       int64  `json:"createTime,omitempty"`
	ModifyTime       int64  `json:"modifyTime,omitempty"`
	Desc             string `json:"desc,omitempty"`
	Tags             string `json:"configTags,omitempty"`
}

func (c *Config) WriteJson(w io.Writer) error {
	return writeJson(c, w)
}

func (c *Config) WriteYaml(w io.Writer) error {
	return writeYaml(c, w)
}

func (c *Config) WriteFile(name string) error {
	return writeYamlFile(c, name)
}

// LoadFromYaml load from YAML file
func (c *Config) LoadFromYaml(name string) error {
	return readYamlFile(c, name)
}

type NsList struct {
	// Code    int          `json:"code,omitempty"`
	// Message interface{}  `json:"message,omitempty"`
	Items []*Namespace `json:"data"`
}

func (n *NsList) WriteTable(w io.Writer) {
	writeTable(w, func(t table.Writer) {
		t.AppendHeader(table.Row{"NAME", "ID", "DESCRIPTION", "COUNT"})
		for _, ns := range n.Items {
			t.AppendRow(table.Row{ns.ShowName, ns.Name, ns.Desc, ns.ConfigCount})
		}
		t.SortBy([]table.SortBy{{Name: "NAME", Mode: table.Asc}, {Name: "ID", Mode: table.Asc}})
	})
}

func (n *NsList) WriteJson(w io.Writer) error {
	return writeJson(n, w)
}

func (n *NsList) WriteYaml(w io.Writer) error {
	return writeYaml(n, w)
}

func (n *NsList) WriteToDir(name string) error {
	for _, c := range n.Items {
		if c.Name == "" {
			continue
		}
		if err := c.WriteFile(filepath.Join(name, fmt.Sprintf("%s.yaml", c.ShowName))); err != nil {
			return err
		}
	}
	return nil
}

type Namespace struct {
	Name        string `json:"namespace"`
	ShowName    string `json:"namespaceShowName"`
	Desc        string `json:"namespaceDesc"`
	Quota       int    `json:"quota"`
	ConfigCount int    `json:"configCount"`
	Type        int    `json:"type"`
}

func (n *Namespace) WriteJson(w io.Writer) error {
	return writeJson(n, w)
}

func (n *Namespace) WriteYaml(w io.Writer) error {
	return writeYaml(n, w)
}
func (n *Namespace) WriteFile(name string) error {
	return writeYamlFile(n, name)
}

func (n *Namespace) LoadFromYaml(name string) error {
	return readYamlFile(n, name)
}
