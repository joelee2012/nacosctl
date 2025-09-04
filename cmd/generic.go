package cmd

import (
	"io"

	"github.com/jedib0t/go-pretty/table"
)

type ItemTypes interface {
	TableRow
	DirWriter
}

type ObjectList[T ItemTypes] struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Items      []T    `json:"items"`
}

func NewList[T ItemTypes, S any](apiVersion string, items []*S, covert func(apiVersion string, s *S) *T) *ObjectList[T] {
	list := new(ObjectList[T])
	list.Kind = "List"
	list.APIVersion = apiVersion
	for _, item := range items {
		list.Items = append(list.Items, *covert(apiVersion, item))
	}
	return list
}

func (o *ObjectList[T]) ToTable(w io.Writer) {
	tb := table.NewWriter()
	tb.SetOutputMirror(w)
	if len(o.Items) == 0 {
		return
	}
	tb.AppendHeader(o.Items[0].TableHeader())
	for _, it := range o.Items {
		tb.AppendRow(it.TableRow())
	}
	tb.SortBy([]table.SortBy{{Name: "NAME", Mode: table.Asc}, {Name: "ID", Mode: table.Asc}})
	s := table.StyleLight
	s.Options = table.OptionsNoBordersAndSeparators
	tb.SetStyle(s)
	tb.Render()
}

func (o *ObjectList[T]) WriteToDir(base string) error {
	for _, it := range o.Items {
		if err := it.WriteToDir(base); err != nil {
			return err
		}
	}
	return nil
}
