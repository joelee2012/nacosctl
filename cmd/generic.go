package cmd

import (
	"io"

	"github.com/jedib0t/go-pretty/table"
	"github.com/joelee2012/nacosctl/pkg/nacos"
)

type ItemTypes interface {
	User | Role | Permission | Configuration | Namespace
	TableWriter
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

func (o *ObjectList[T]) WriteToDir(base string) error {
	for _, it := range o.Items {
		if err := it.WriteToDir(base); err != nil {
			return err
		}
	}
	return nil
}

type FormatWriter interface {
	TableWriter
	DirWriter
}

func WriteFormat[T ItemTypes](obj *ObjectList[T], format string, w io.Writer) error {
	switch format {
	case "json":
		return toJson(obj, w)
	case "yaml":
		return toYaml(obj, w)
	case "table":
		obj.ToTable(w)
	default:
		return obj.WriteToDir(format)
	}
	return nil
}
