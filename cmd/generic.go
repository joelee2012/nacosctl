package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jedib0t/go-pretty/table"
	"github.com/joelee2012/nacosctl/pkg/nacos"
)

type TableRow interface {
	TableHeader() table.Row
	TableRow() table.Row
}

type ItemTypes interface {
	User | Role | Permission | Configuration | Namespace
	TableRow
	DirWriter
}

type SrcTypes interface {
	nacos.User | nacos.Role | nacos.Permission | nacos.Configuration | nacos.Namespace
}

type ObjectList[T ItemTypes] struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Items      []T    `json:"items"`
}

func NewList[T ItemTypes, S SrcTypes](apiVersion string, items []*S, newFunc func(apiVersion string, s *S) *T) *ObjectList[T] {
	list := new(ObjectList[T])
	list.Kind = "List"
	list.APIVersion = apiVersion
	for _, item := range items {
		list.Items = append(list.Items, *newFunc(apiVersion, item))
	}
	return list
}

func (o *ObjectList[T]) ToTable(w io.Writer) {
	toTable(w, func(t table.Writer) {
		t.AppendHeader(o.Items[0].TableHeader())
		for _, it := range o.Items {
			t.AppendRow(it.TableRow())
		}
		t.SortBy([]table.SortBy{{Name: "NAME", Mode: table.Asc}, {Name: "ID", Mode: table.Asc}})
	})
}

func PrintTable[T TableRow](items []T, w io.Writer) {
	tw := table.NewWriter()
	tw.SetOutputMirror(w)
	if len(items) == 0 {
		return
	}
	tw.AppendHeader(items[0].TableHeader())
	for _, it := range items {
		tw.AppendRow(it.TableRow())
	}
	s := table.StyleLight
	s.Options = table.OptionsNoBordersAndSeparators
	tw.SetStyle(s)
	tw.Render()
}

// 把列表里的每个元素写入目录
func WriteDir[T DirWriter](items []T, base string) error {
	for _, it := range items {
		if err := it.WriteToDir(base); err != nil {
			return err
		}
	}
	return nil
}

func (c Configuration) TableHeader() table.Row {
	return table.Row{"NAMESPACEID", "DATAID", "GROUP", "APPLICATION", "TYPE"}
}

func (c Configuration) TableRow() table.Row {
	return table.Row{c.Metadata.Namespace, c.Metadata.DataID, c.Metadata.Group,
		c.Spec.Application, c.Spec.Type}
}

func (c Configuration) WriteToDir(base string) error {
	dir := filepath.Join(base, c.Metadata.Namespace, c.Metadata.Group)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}
	return c.ToFile(filepath.Join(dir, c.Metadata.DataID))
}

func (n Namespace) TableHeader() table.Row {
	return table.Row{"NAME", "ID", "DESCRIPTION", "COUNT"}
}
func (n Namespace) TableRow() table.Row {
	return table.Row{n.Metadata.Name, n.Metadata.ID, n.Metadata.Description,
		fmt.Sprintf("%d", n.Status.ConfigCount)}
}
func (n Namespace) WriteToDir(base string) error {
	if n.Metadata.ID == "" {
		return nil
	}
	return n.ToFile(filepath.Join(base, n.Metadata.ID+".yaml"))
}

func WriteFormat[T ItemTypes](obj *ObjectList[T], format string, w io.Writer) error {
	switch format {
	case "json":
		return toJson(obj, w)
	case "yaml":
		return toYaml(obj, w)
	case "table":
		PrintTable(obj.Items, w)
	default:
		return WriteDir(obj.Items, format)
	}
	return nil
}
