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

type YamlFileLoader interface {
	LoadFromYaml(name string) error
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

func (c *Config) WriteJson(w io.Writer) error {
	return writeJson(c, w)
}

func (c *Config) WriteYaml(w io.Writer) error {
	return writeYaml(c, w)
}

func (c *Config) WriteFile(name string) error {
	return writeYamlFile(c, name)
}

func (c *Config) LoadFromYaml(name string) error {
	return readYamlFile(c, name)
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
