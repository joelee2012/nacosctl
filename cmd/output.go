package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
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

func writeJson(v any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func writeYaml(v any, w io.Writer) error {
	enc := yaml.NewEncoder(w)
	return enc.Encode(v)
}

func writeFile(v any, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return writeYaml(v, f)
}

func writeTable(w io.Writer, fn func(t table.Writer)) {
	tb := table.NewWriter()
	tb.SetOutputMirror(w)
	fn(tb)
	s := table.StyleLight
	s.Options = table.OptionsNoBordersAndSeparators
	tb.SetStyle(s)
	tb.Render()
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
	return writeFile(c, name)
}

func readYaml(v any, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(v)
}

func (c *Config) LoadFromYaml(name string) error {
	return readYaml(c, name)
}

func (n *NsList) WriteTable(w io.Writer) {
	writeTable(w, func(t table.Writer) {
		t.AppendHeader(table.Row{"NAMESPACE", "ID", "DESCRIPTION", "COUNT"})
		for _, ns := range n.Items {
			t.AppendRow(table.Row{ns.ShowName, ns.Name, ns.Desc, ns.ConfigCount})
		}
		t.SortBy([]table.SortBy{{Name: "NAMESPACE", Mode: table.Asc}, {Name: "ID", Mode: table.Asc}})
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
		if err := c.WriteFile(filepath.Join(cmdOpts.OutDir, fmt.Sprintf("%s.yaml", c.ShowName))); err != nil {
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
	return writeFile(n, name)
}

func (n *Namespace) LoadFromYaml(name string) error {
	return readYaml(n, name)
}
func WriteAsFormat(format string, writable FormatWriter) {
	switch format {
	case "json":
		writable.WriteJson(os.Stdout)
	case "yaml":
		writable.WriteYaml(os.Stdout)
	case "table":
		writable.WriteTable(os.Stdout)
	default:
		writable.WriteTable(os.Stdout)
	}
}

func LoadFromYaml(name string, loader YamlFileLoader) error {
	return loader.LoadFromYaml(name)
}
