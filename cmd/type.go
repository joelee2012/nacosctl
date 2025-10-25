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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/jedib0t/go-pretty/table"
	"github.com/joelee2012/nacosctl/pkg/nacos"
)

type DirWriter interface {
	WriteToDir(name string) error
}
type TableWriter interface {
	ToTable(w io.Writer)
}

type TableRow interface {
	TableHeader() table.Row
	TableRow() table.Row
}

type Configuration struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Group     string `json:"group"`
		DataID    string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Content     string `json:"data"`
		Type        string `json:"type"`
		Application string `json:"application,omitempty"`
		Description string `json:"description,omitempty"`
		Tags        string `json:"tags,omitempty"`
	} `json:"spec"`
	Status struct {
		Md5              string `json:"md5,omitempty"`
		EncryptedDataKey string `json:"encryptedDataKey,omitempty"`
		CreateTime       int64  `json:"createTime,omitempty"`
		ModifyTime       int64  `json:"modifyTime,omitempty"`
	} `json:"status"`
}

func NewConfiguration(apiVersion string, cfg nacos.Configuration) *Configuration {
	c := new(Configuration)
	c.APIVersion = apiVersion
	c.Kind = "Configuration"
	c.Metadata.DataID = cfg.DataID
	c.Metadata.Namespace = cfg.GetNamespace()
	c.Metadata.Group = cfg.GetGroup()
	c.Spec.Application = cfg.Application
	c.Spec.Type = cfg.Type
	c.Spec.Tags = cfg.Tags
	c.Spec.Description = cfg.Description
	c.Spec.Content = cfg.Content
	c.Status.Md5 = cfg.Md5
	c.Status.CreateTime = cfg.CreateTime
	c.Status.ModifyTime = cfg.ModifyTime
	c.Status.EncryptedDataKey = cfg.EncryptedDataKey
	return c
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
	return writeYamlFile(c, filepath.Join(dir, c.Metadata.DataID))
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
		Quota       int `json:"quota,omitempty"`
		ConfigCount int `json:"configCount,omitempty"`
		Type        int `json:"type,omitempty"`
	} `json:"status"`
}

func NewNamespace(apiVersion string, ns nacos.Namespace) *Namespace {
	n := new(Namespace)
	n.APIVersion = apiVersion
	n.Kind = "Namespace"
	n.Metadata.Name = ns.Name
	n.Metadata.ID = ns.ID
	n.Metadata.Description = ns.Description
	n.Status.ConfigCount = ns.ConfigCount
	n.Status.Quota = ns.Quota
	n.Status.Type = ns.Type
	return n
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
	return writeYamlFile(n, filepath.Join(base, n.Metadata.ID+".yaml"))
}

type User struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	} `json:"metadata"`
}

func NewUser(apiVersion string, user nacos.User) *User {
	u := new(User)
	u.APIVersion = apiVersion
	u.Kind = "User"
	u.Metadata.Name = user.Name
	u.Metadata.Password = user.Password
	return u
}

func (u User) TableHeader() table.Row {
	return table.Row{"NAME", "PASSWORD"}
}
func (u User) TableRow() table.Row {
	return table.Row{u.Metadata.Name, u.Metadata.Password}
}
func (u User) WriteToDir(base string) error {
	return writeYamlFile(u, filepath.Join(base, u.Metadata.Name+".yaml"))
}

type Role struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"metadata"`
}

func NewRole(apiVersion string, role nacos.Role) *Role {
	r := new(Role)
	r.APIVersion = apiVersion
	r.Kind = "Role"
	r.Metadata.Name = role.Name
	r.Metadata.Username = role.Username
	return r
}

func (r Role) TableHeader() table.Row {
	return table.Row{"NAME", "USERNAME"}
}
func (r Role) TableRow() table.Row {
	return table.Row{r.Metadata.Name, r.Metadata.Username}
}
func (r Role) WriteToDir(base string) error {
	return writeYamlFile(r, filepath.Join(base, r.Metadata.Name+r.Metadata.Username+".yaml"))
}

type Permission struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Role     string `json:"role"`
		Resource string `json:"resource"`
		Action   string `json:"action"`
	} `json:"metadata"`
}

func NewPermission(apiVersion string, perm nacos.Permission) *Permission {
	p := new(Permission)
	p.APIVersion = apiVersion
	p.Kind = "Permission"
	p.Metadata.Role = perm.Role
	p.Metadata.Resource = perm.Resource
	p.Metadata.Action = perm.Action
	return p
}

func (p Permission) TableHeader() table.Row {
	return table.Row{"ROLE", "RESOURCE", "ACTION"}
}
func (p Permission) TableRow() table.Row {
	return table.Row{p.Metadata.Role, p.Metadata.Resource, p.Metadata.Action}
}
func (p Permission) WriteToDir(base string) error {
	return writeYamlFile(p, filepath.Join(base, p.Metadata.Role+p.Metadata.Resource+p.Metadata.Action+".yaml"))
}

func toJson(v any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func toYaml(v any, w io.Writer) error {
	enc := yaml.NewEncoder(w, yaml.UseLiteralStyleIfMultiline(true))
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

type FormatWriter interface {
	TableWriter
	DirWriter
}

func WriteFormat(fw FormatWriter, format string, w io.Writer) error {
	switch format {
	case "json":
		return toJson(fw, w)
	case "yaml":
		return toYaml(fw, w)
	case "table":
		fw.ToTable(w)
	default:
		return fw.WriteToDir(format)
	}
	return nil
}
