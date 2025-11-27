/*
Copyright © 2025 Joe Lee <lj_2005@163.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"io"

	"github.com/jedib0t/go-pretty/table"
)

type ListTypes interface {
	TableRow
	DirWriter
}

type List[T ListTypes] struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Items      []T    `json:"items"`
}

type ConfigurationList = List[Configuration]
type NamespaceList = List[Namespace]
type PermissionList = List[Permission]
type RoleList = List[Role]
type UserList = List[User]

func NewList[T ListTypes, S any](apiVersion string, items []S, covert func(apiVersion string, s S) *T) *List[T] {
	list := new(List[T])
	list.Kind = "List"
	list.APIVersion = apiVersion
	for _, item := range items {
		list.Items = append(list.Items, *covert(apiVersion, item))
	}
	return list
}

func (lst *List[T]) ToTable(w io.Writer) {
	tb := table.NewWriter()
	tb.SetOutputMirror(w)
	if len(lst.Items) == 0 {
		w.Write([]byte("No resources found"))
		return
	}
	tb.AppendHeader(lst.Items[0].TableHeader())
	for _, it := range lst.Items {
		tb.AppendRow(it.TableRow())
	}
	tb.SortBy([]table.SortBy{{Name: "NAME", Mode: table.Asc}, {Name: "ID", Mode: table.Asc}})
	s := table.StyleLight
	s.Options = table.OptionsNoBordersAndSeparators
	tb.SetStyle(s)
	tb.Render()
}

func (lst *List[T]) WriteToDir(base string) error {
	for _, it := range lst.Items {
		if err := it.WriteToDir(base); err != nil {
			return err
		}
	}
	return nil
}
