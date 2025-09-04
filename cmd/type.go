package cmd

import (
	"encoding/json"
	"io"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/jedib0t/go-pretty/table"
	"github.com/joelee2012/nacosctl/pkg/nacos"
)

type DirWriter interface {
	WriteToDir(name string) error
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

type User struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	} `json:"metadata"`
}

func NewUser(apiVersion string, user *nacos.User) *User {
	u := new(User)
	u.APIVersion = apiVersion
	u.Kind = "User"
	u.Metadata.Name = user.Name
	u.Metadata.Password = user.Password
	return u
}

type Role struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"metadata"`
}

func NewRole(apiVersion string, user *nacos.Role) *Role {
	r := new(Role)
	r.APIVersion = apiVersion
	r.Kind = "Role"
	r.Metadata.Name = user.Name
	r.Metadata.Username = user.Username
	return r
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

func NewPermission(apiVersion string, perm *nacos.Permission) *Permission {
	p := new(Permission)
	p.APIVersion = apiVersion
	p.Kind = "Permission"
	p.Metadata.Role = perm.Role
	p.Metadata.Resource = perm.Resource
	p.Metadata.Action = perm.Action
	return p
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

func toTable(w io.Writer, fn func(t table.Writer)) {
	tb := table.NewWriter()
	tb.SetOutputMirror(w)
	fn(tb)
	s := table.StyleLight
	s.Options = table.OptionsNoBordersAndSeparators
	tb.SetStyle(s)
	tb.Render()
}
