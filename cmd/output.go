package cmd

import (
	"encoding/json"
	"io"
	"os"

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

func writeJson(v any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func writeYaml(v any, w io.Writer) error {
	enc := yaml.NewEncoder(w)
	// enc.SetIndent(2)
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

func (c *ConfigList) WriteTable(w io.Writer) {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.AppendHeader(table.Row{"NAMESPACE", "DATAID", "GROUP", "APPLICATION", "TYPE"})
	for _, item := range c.PageItems {
		t.AppendRow(table.Row{item.Tenant, item.DataID, item.Group, item.AppName, item.Type})
	}
	t.SortBy([]table.SortBy{{Name: "NAMESPACE", Mode: table.Asc}, {Name: "DATAID", Mode: table.Asc}})
	s := table.StyleLight
	s.Options = table.OptionsNoBordersAndSeparators
	t.SetStyle(s)
	t.Render()
}

func (c *ConfigList) WriteJson(w io.Writer) error {
	return writeJson(c, w)
}

func (c *ConfigList) WriteYaml(w io.Writer) error {
	return writeYaml(c, w)
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

func (n *NsList) WriteTable(w io.Writer) {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.AppendHeader(table.Row{"NAMESPACE", "ID", "DESCRIPTION", "COUNT"})
	for _, ns := range n.Items {
		t.AppendRow(table.Row{ns.ShowName, ns.Name, ns.Desc, ns.ConfigCount})
	}
	t.SortBy([]table.SortBy{{Name: "NAMESPACE", Mode: table.Asc}, {Name: "ID", Mode: table.Asc}})
	s := table.StyleLight
	s.Options = table.OptionsNoBordersAndSeparators
	t.SetStyle(s)
	t.Render()
}

func (n *NsList) WriteJson(w io.Writer) error {
	return writeJson(n, w)
}

func (n *NsList) WriteYaml(w io.Writer) error {
	return writeYaml(n, w)
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
