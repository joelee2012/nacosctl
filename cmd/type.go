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
	DirWriter
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

func NewConfiguration(apiVersion string, nc *nacos.Configuration) *Configuration {
	c := new(Configuration)
	c.APIVersion = apiVersion
	c.Kind = "Configuration"
	c.Metadata.DataID = nc.DataID
	c.Metadata.Namespace = nc.NamespaceID
	c.Metadata.Group = nc.Group
	c.Spec.Application = nc.Application
	c.Spec.Type = nc.Type
	c.Spec.Tags = nc.Tags
	c.Spec.Description = nc.Description
	c.Spec.Content = nc.Content
	c.Status.Md5 = nc.Md5
	c.Status.CreateTime = nc.CreateTime
	c.Status.ModifyTime = nc.ModifyTime
	c.Status.EncryptedDataKey = nc.EncryptedDataKey
	return c
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

type ConfigurationList struct {
	APIVersion string           `json:"apiVersion"`
	Items      []*Configuration `json:"items"`
	Kind       string           `json:"kind"`
}

func NewConfigurationList(apiVersion string, cs *nacos.ConfigurationList) *ConfigurationList {
	list := new(ConfigurationList)
	list.Kind = "List"
	list.APIVersion = apiVersion
	for _, item := range cs.Items {
		list.Items = append(list.Items, NewConfiguration(apiVersion, item))
	}
	return list
}

func (list *ConfigurationList) ToTable(w io.Writer) {
	toTable(w, func(t table.Writer) {
		t.AppendHeader(table.Row{"NAMESPACEID", "DATAID", "GROUP", "APPLICATION", "TYPE"})
		for _, item := range list.Items {
			t.AppendRow(table.Row{item.Metadata.Namespace, item.Metadata.DataID, item.Metadata.Group, item.Spec.Application, item.Spec.Type})
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

func NewNamespaceList(apiVersion string, ns *nacos.NamespaceList) *NamespaceList {
	list := new(NamespaceList)
	list.Kind = "List"
	list.APIVersion = apiVersion
	for _, item := range ns.Items {
		list.Items = append(list.Items, NewNamespace(apiVersion, item))
	}
	return list
}

func (list *NamespaceList) ToJson(w io.Writer) error {
	return toJson(list, w)
}

func (list *NamespaceList) ToYaml(w io.Writer) error {
	return toYaml(list, w)
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
		Quota       int `json:"quota,omitempty"`
		ConfigCount int `json:"configCount,omitempty"`
		Type        int `json:"type,omitempty"`
	} `json:"status"`
}

func NewNamespace(apiVersion string, e *nacos.Namespace) *Namespace {
	n := new(Namespace)
	n.APIVersion = apiVersion
	n.Kind = "Namespace"
	n.Metadata.Name = e.Name
	n.Metadata.ID = e.ID
	n.Metadata.Description = e.Description
	n.Status.ConfigCount = e.ConfigCount
	n.Status.Quota = e.Quota
	n.Status.Type = e.Type
	return n
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

func toJson(v any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func toYaml(v any, w io.Writer) error {
	enc := yaml.NewEncoder(w, yaml.UseLiteralStyleIfMultiline(true))
	return enc.Encode(v)
}

type UserList struct {
	APIVersion string  `json:"apiVersion"`
	Kind       string  `json:"kind"`
	Items      []*User `json:"items"`
}

type User struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	} `json:"metadata"`
}

func NewUserList(apiVersion string, users *nacos.UserList) *UserList {
	list := new(UserList)
	list.Kind = "List"
	list.APIVersion = apiVersion
	for _, item := range users.Items {
		list.Items = append(list.Items, NewUser(apiVersion, item))
	}
	return list
}

func NewUser(apiVersion string, user *nacos.User) *User {
	u := new(User)
	u.APIVersion = apiVersion
	u.Kind = "User"
	u.Metadata.Name = user.Name
	u.Metadata.Password = user.Password
	return u
}

// func (list *UserList) ToJson(w io.Writer) error {
// 	return toJson(list, w)
// }

// func (list *UserList) ToYaml(w io.Writer) error {
// 	return toYaml(list, w)
// }

// func (list *UserList) ToTable(w io.Writer) {
// 	toTable(w, func(t table.Writer) {
// 		t.AppendHeader(table.Row{"USERNAME"})
// 		for _, ns := range list.Items {
// 			t.AppendRow(table.Row{ns.Metadata.Name, ns.Metadata.ID, ns.Metadata.Description, ns.Status.ConfigCount})
// 		}
// 		t.SortBy([]table.SortBy{{Name: "NAME", Mode: table.Asc}, {Name: "ID", Mode: table.Asc}})
// 	})
// }

// func (n *NamespaceList) WriteToDir(name string) error {
// 	for _, e := range n.Items {
// 		if e.Metadata.ID == "" {
// 			continue
// 		}
// 		if err := e.ToFile(filepath.Join(name, fmt.Sprintf("%s.yaml", e.Metadata.ID))); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

type RoleList struct {
	APIVersion string  `json:"apiVersion"`
	Kind       string  `json:"kind"`
	Items      []*Role `json:"items"`
}

type Role struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"metadata"`
}

func NewRoleList(apiVersion string, roles *nacos.RoleList) *RoleList {
	list := new(RoleList)
	list.Kind = "List"
	list.APIVersion = apiVersion
	for _, item := range roles.Items {
		list.Items = append(list.Items, NewRole(apiVersion, item))
	}
	return list
}

func NewRole(apiVersion string, user *nacos.Role) *Role {
	r := new(Role)
	r.APIVersion = apiVersion
	r.Kind = "Role"
	r.Metadata.Name = user.Name
	r.Metadata.Username = user.Username
	return r
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

func WriteAsFormat(format string, writable FormatWriter, w io.Writer) error {
	switch format {
	case "json":
		return writable.ToJson(w)
	case "yaml":
		return writable.ToYaml(w)
	case "table":
		writable.ToTable(w)
	default:
		return writable.WriteToDir(format)
	}
	return nil
}

func WriteToDir(name string, writable DirWriter) error {
	return writable.WriteToDir(name)
}
func LoadFromYaml(name string, loader YamlFileLoader) error {
	return loader.FromYaml(name)
}
