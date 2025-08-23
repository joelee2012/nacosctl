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

type ConfigList struct {
	TotalCount     int       `json:"totalCount,omitempty"`
	PageNumber     int       `json:"pageNumber,omitempty"`
	PagesAvailable int       `json:"pagesAvailable,omitempty"`
	Items          []*Config `json:"pageItems"`
}

func (c *ConfigList) ToTable(w io.Writer) {
	toTable(w, func(t table.Writer) {
		t.AppendHeader(table.Row{"NAMESPACEID", "DATAID", "GROUP", "APPLICATION", "TYPE"})
		for _, item := range c.Items {
			t.AppendRow(table.Row{item.NamespaceID, item.DataID, item.Group, item.Application, item.Type})
		}
		t.SortBy([]table.SortBy{{Name: "NAMESPACEID", Mode: table.Asc}, {Name: "DATAID", Mode: table.Asc}})
	})
}

func (c *ConfigList) ToJson(w io.Writer) error {
	return toJson(c, w)
}

func (c *ConfigList) ToYaml(w io.Writer) error {
	return toYaml(c, w)
}

func (cs *ConfigList) WriteToDir(name string) error {
	var dir string
	for _, c := range cs.Items {
		if c.NamespaceID == "" {
			dir = filepath.Join(name, "public", c.Group)
		} else {
			dir = filepath.Join(name, c.NamespaceID, c.Group)
		}
		if err := os.MkdirAll(dir, 0750); err != nil {
			return err
		}
		if err := c.ToFile(filepath.Join(dir, c.DataID)); err != nil {
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
	NamespaceID      string `json:"tenant"`
	Type             string `json:"type"`
	Md5              string `json:"md5,omitempty"`
	EncryptedDataKey string `json:"encryptedDataKey,omitempty"`
	Application      string `json:"appName,omitempty"`
	CreateTime       int64  `json:"createTime,omitempty"`
	ModifyTime       int64  `json:"modifyTime,omitempty"`
	Description      string `json:"desc,omitempty"`
	Tags             string `json:"configTags,omitempty"`
}

func (c *Config) ToJson(w io.Writer) error {
	return toJson(c, w)
}

func (c *Config) ToYaml(w io.Writer) error {
	return toYaml(c, w)
}

func (c *Config) ToFile(name string) error {
	return writeYamlFile(c, name)
}

// FromYaml load from YAML file
func (c *Config) FromYaml(name string) error {
	return readYamlFile(c, name)
}

type NsList struct {
	// Code    int          `json:"code,omitempty"`
	// Message interface{}  `json:"message,omitempty"`
	Items []*Namespace `json:"data"`
}

func (n *NsList) ToTable(w io.Writer) {
	toTable(w, func(t table.Writer) {
		t.AppendHeader(table.Row{"NAME", "ID", "DESCRIPTION", "COUNT"})
		for _, ns := range n.Items {
			t.AppendRow(table.Row{ns.Name, ns.ID, ns.Description, ns.ConfigCount})
		}
		t.SortBy([]table.SortBy{{Name: "NAME", Mode: table.Asc}, {Name: "ID", Mode: table.Asc}})
	})
}

func (n *NsList) ToJson(w io.Writer) error {
	return toJson(n, w)
}

func (n *NsList) ToYaml(w io.Writer) error {
	return toYaml(n, w)
}

func (n *NsList) WriteToDir(name string) error {
	for _, e := range n.Items {
		if e.ID == "" {
			continue
		}
		if err := e.ToFile(filepath.Join(name, fmt.Sprintf("%s.yaml", e.Name))); err != nil {
			return err
		}
	}
	return nil
}

type Namespace struct {
	ID          string `json:"namespace"`
	Name        string `json:"namespaceShowName"`
	Description string `json:"namespaceDesc"`
	Quota       int    `json:"quota"`
	ConfigCount int    `json:"configCount"`
	Type        int    `json:"type"`
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

type UserList struct {
	TotalCount     int     `json:"totalCount,omitempty"`
	PageNumber     int     `json:"pageNumber,omitempty"`
	PagesAvailable int     `json:"pagesAvailable,omitempty"`
	Items          []*User `json:"pageItems"`
}

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type RoleList struct {
	TotalCount     int     `json:"totalCount,omitempty"`
	PageNumber     int     `json:"pageNumber,omitempty"`
	PagesAvailable int     `json:"pagesAvailable,omitempty"`
	Items          []*Role `json:"pageItems"`
}

type Role struct {
	Name     string `json:"role"`
	Username string `json:"username"`
}

type PermissionList struct {
	TotalCount     int           `json:"totalCount,omitempty"`
	PageNumber     int           `json:"pageNumber,omitempty"`
	PagesAvailable int           `json:"pagesAvailable,omitempty"`
	Items          []*Permission `json:"pageItems"`
}

type Permission struct {
	Role       string `json:"role"`
	Resource   string `json:"resource"`
	Permission string `json:"action"`
}
