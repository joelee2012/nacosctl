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
package nacos

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	URL        string
	User       string
	Password   string
	APIVersion string
	*Token
	*State
}
type Token struct {
	AccessToken string `json:"accessToken"`
	TokenTTL    int64  `json:"tokenTtl"`
	GlobalAdmin bool   `json:"globalAdmin"`
	Username    string `json:"username"`
	ExpiredAt   int64
}

func (t *Token) Expired() bool {
	return time.Now().After(time.Unix(t.ExpiredAt, 0))
}

type State struct {
	Version        string `json:"version"`
	StandaloneMode string `json:"standalone_mode"`
	FunctionMode   string `json:"function_mode"`
}

var api = map[string]map[string]string{
	"v1": {
		"state":     "/v1/console/server/state",
		"token":     "/v1/auth/login",
		"list_ns":   "/v1/console/namespaces",
		"ns":        "/v1/console/namespaces",
		"cs":        "/v1/cs/configs",
		"list_cs":   "/v1/cs/configs",
		"user":      "/v1/auth/users",
		"list_user": "/v1/auth/users",
		"role":      "/v1/auth/roles",
		"list_role": "/v1/auth/roles",
		"perm":      "/v1/auth/permissions",
		"list_perm": "/v1/auth/permissions",
	},
	"v3": {
		"state":     "/v3/console/server/state",
		"token":     "/v3/auth/user/login",
		"list_ns":   "/v3/console/core/namespace/list",
		"ns":        "/v3/console/core/namespace",
		"cs":        "/v3/console/cs/config",
		"list_cs":   "/v3/console/cs/config/list",
		"list_user": "/v3/auth/user/list",
		"user":      "/v3/auth/user",
		"list_role": "/v3/auth/role/list",
		"role":      "/v3/auth/role",
		"perm":      "/v3/auth/permission",
		"list_perm": "/v3/auth/permission/list",
	},
}

var apiVersion = []string{"v3", "v1"}

func NewClient(url, user, password string) *Client {
	return &Client{
		URL:        url,
		User:       user,
		Password:   password,
		APIVersion: "v1",
	}
}

func (c *Client) DetectAPIVersion() {
	for _, ver := range apiVersion {
		c.APIVersion = ver
		v, err := c.GetVersion()
		if err == nil && v != "" {
			return
		}
	}
	c.APIVersion = "v1"
}

func (c *Client) GetVersion() (string, error) {
	if c.State != nil {
		return c.Version, nil
	}
	resp, err := http.Get(c.URL + api[c.APIVersion]["state"])
	err = decode(resp, err, &c.State)
	if err != nil {
		return "", err
	}
	return c.Version, err
}

func (c *Client) GetToken() (string, error) {
	if c.Token != nil && !c.Token.Expired() {
		return c.AccessToken, nil
	}
	v := url.Values{}
	v.Add("username", c.User)
	v.Add("password", c.Password)
	now := time.Now().Unix()
	resp, err := http.PostForm(c.URL+api[c.APIVersion]["token"], v)
	err = decode(resp, err, &c.Token)
	if err != nil {
		return "", err
	}
	c.ExpiredAt = now + c.TokenTTL
	return c.AccessToken, err
}

func (c *Client) ListNamespace() (*NamespaceList, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s%s?%s", c.URL, api[c.APIVersion]["list_ns"], v.Encode())
	resp, err := http.Get(url)
	namespaces := new(NamespaceList)
	err = decode(resp, err, namespaces)
	return namespaces, err
}

type CreateNsOpts struct {
	Name        string
	Description string
	ID          string
}

func (c *Client) CreateNamespace(opts *CreateNsOpts) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("customNamespaceId", opts.ID)
	v.Add("namespaceName", opts.Name)
	v.Add("namespaceDesc", opts.Description)
	v.Add("accessToken", token)
	resp, err := http.PostForm(c.URL+api[c.APIVersion]["ns"], v)
	return checkErr(resp, err)
}

func (c *Client) DeleteNamespace(id string) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("namespaceId", id)
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s%s?%s", c.URL, api[c.APIVersion]["ns"], v.Encode())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	return checkErr(resp, err)
}

func (c *Client) UpdateNamespace(opts *CreateNsOpts) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("namespace", opts.ID)
	v.Add("namespaceShowName", opts.Name)
	v.Add("namespaceDesc", opts.Description)
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s%s?%s", c.URL, api[c.APIVersion]["ns"], v.Encode())
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	return checkErr(resp, err)
}

func (c *Client) CreateOrUpdateNamespace(opts *CreateNsOpts) error {
	nsList, err := c.ListNamespace()
	if err != nil {
		return err
	}
	for _, ns := range nsList.Items {
		if ns.ID == opts.ID {
			return c.UpdateNamespace(opts)
		}
	}
	return c.CreateNamespace(opts)
}

func (c *Client) GetNamespace(id string) (*Namespace, error) {
	nsList, err := c.ListNamespace()
	if err != nil {
		return nil, err
	}
	for _, ns := range nsList.Items {
		if ns.ID == id {
			return &ns, nil
		}
	}
	return nil, fmt.Errorf("404 Not Found %s", id)
}

type GetCfgOpts struct {
	DataID      string
	Group       string
	NamespaceID string
}

func (c *Client) GetConfig(opts *GetCfgOpts) (*Configuration, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Add("dataId", opts.DataID)
	v.Add("group", opts.Group)
	v.Add("groupName", opts.Group)
	v.Add("namespaceId", opts.NamespaceID)
	v.Add("tenant", opts.NamespaceID)
	v.Add("show", "all")
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s%s?%s", c.URL, api[c.APIVersion]["cs"], v.Encode())
	resp, err := http.Get(url)
	if c.APIVersion == "v3" {
		config := new(ConfigurationV3)
		err = decode(resp, err, config)
		if err == io.EOF {
			return nil, fmt.Errorf("404 Not Found %s %w", url, err)
		}
		return &config.Data, nil
	} else {

		config := new(Configuration)
		err = decode(resp, err, config)
		// if config not found, nacos server return 200 and empty response
		if err == io.EOF {
			return nil, fmt.Errorf("404 Not Found %s %w", url, err)
		}
		return config, err
	}
}

type ListCfgOpts struct {
	Application string
	Content     string
	DataID      string
	Group       string
	NamespaceID string
	Tags        string
	PageNumber  int
	PageSize    int
}

func (c *Client) ListConfig(opts *ListCfgOpts) (*ConfigurationList, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Add("dataId", opts.DataID)
	v.Add("group", opts.Group)
	v.Add("groupName", opts.Group)
	v.Add("appName", opts.Application)
	v.Add("config_tags", opts.Tags)
	v.Add("configTags", opts.Tags)
	if opts.PageNumber == 0 {
		opts.PageNumber = 1
	}
	if opts.PageSize == 0 {
		opts.PageSize = 10
	}
	v.Add("pageNo", strconv.Itoa(opts.PageNumber))
	v.Add("pageSize", strconv.Itoa(opts.PageSize))
	v.Add("tenant", opts.NamespaceID)
	v.Add("namespaceId", opts.NamespaceID)
	v.Add("search", "accurate")
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s%s?%s", c.URL, api[c.APIVersion]["list_cs"], v.Encode())
	if c.APIVersion == "v1" {
		configs := new(ConfigurationList)
		resp, err := http.Get(url)
		err = decode(resp, err, configs)
		return configs, err
	} else {
		configs := new(ConfigurationListV3)
		resp, err := http.Get(url)
		err = decode(resp, err, configs)
		return configs.Data, err
	}
}

func (c *Client) ListConfigInNs(namespace, group string) (*ConfigurationList, error) {
	nsCs := new(ConfigurationList)
	listOpts := ListCfgOpts{PageNumber: 1, PageSize: 100, Group: group, NamespaceID: namespace}
	for {
		cs, err := c.ListConfig(&listOpts)
		if err != nil {
			log.Fatal(err)
		}
		nsCs.Items = append(nsCs.Items, cs.Items...)
		if cs.PagesAvailable == 0 || cs.PagesAvailable == cs.PageNumber {
			break
		}
		listOpts.PageNumber += 1
	}
	return nsCs, nil
}

func (c *Client) ListAllConfig() (*ConfigurationList, error) {
	allCs := new(ConfigurationList)
	nss, err := c.ListNamespace()
	if err != nil {
		return nil, err
	}
	for _, ns := range nss.Items {
		cs, err := c.ListConfigInNs(ns.ID, "")
		if err != nil {
			return nil, err
		}
		allCs.Items = append(allCs.Items, cs.Items...)
	}
	return allCs, nil
}

type CreateCfgOpts struct {
	Application string
	Content     string
	DataID      string
	Description string
	Group       string
	NamespaceID string
	Tags        string
	Type        string
}

func (c *Client) CreateConfig(opts *CreateCfgOpts) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("dataId", opts.DataID)
	v.Add("group", opts.Group)
	v.Add("groupName", opts.Group)
	v.Add("content", opts.Content)
	v.Add("type", opts.Type)
	v.Add("tenant", opts.NamespaceID)
	v.Add("namespaceId", opts.NamespaceID)
	v.Add("appName", opts.Application)
	v.Add("desc", opts.Description)
	v.Add("config_tags", opts.Tags)
	v.Add("accessToken", token)
	resp, err := http.PostForm(c.URL+api[c.APIVersion]["cs"], v)
	return checkErr(resp, err)
}

type DeleteCfgOpts = GetCfgOpts

func (c *Client) DeleteConfig(opts *DeleteCfgOpts) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("dataId", opts.DataID)
	v.Add("group", opts.Group)
	v.Add("groupName", opts.Group)
	v.Add("tenant", opts.NamespaceID)
	v.Add("namespaceId", opts.NamespaceID)
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s%s?%s", c.URL, api[c.APIVersion]["cs"], v.Encode())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	return checkErr(resp, err)
}

func (c *Client) CreateUser(name, password string) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("username", name)
	v.Add("password", password)
	v.Add("accessToken", token)
	resp, err := http.PostForm(c.URL+api[c.APIVersion]["user"], v)
	return checkErr(resp, err)
}

func (c *Client) DeleteUser(name string) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("username", name)
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s%s?%s", c.URL, api[c.APIVersion]["user"], v.Encode())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	return checkErr(resp, err)
}

func (c *Client) ListUser() (*UserList, error) {
	if c.APIVersion == "v1" {
		return listResource[UserList](c, api[c.APIVersion]["list_user"])
	}
	return listResource[UserListV3](c, api[c.APIVersion]["list_user"])
}

func (c *Client) GetUser(name string) (*User, error) {
	users, err := c.ListUser()
	if err != nil {
		return nil, err
	}

	for _, user := range users.Items {
		if user.Name == name {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("404 Not Found %s", name)
}

func (c *Client) CreateRole(name, username string) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("username", username)
	v.Add("role", name)
	v.Add("accessToken", token)
	resp, err := http.PostForm(c.URL+api[c.APIVersion]["role"], v)
	return checkErr(resp, err)
}

func (c *Client) DeleteRole(name, username string) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("username", username)
	v.Add("role", name)
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s%s?%s", c.URL, api[c.APIVersion]["role"], v.Encode())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	return checkErr(resp, err)
}

func (c *Client) ListRole() (*RoleList, error) {
	if c.APIVersion == "v1" {
		return listResource[RoleList](c, api[c.APIVersion]["list_role"])
	}
	return listResource[RoleListV3](c, api[c.APIVersion]["list_role"])
}

func (c *Client) GetRole(name, username string) (*Role, error) {
	roles, err := c.ListRole()
	if err != nil {
		return nil, err
	}
	r := Role{Name: name, Username: username}
	if roles.Contains(r) {
		return &r, nil
	}
	return nil, fmt.Errorf("404 Not Found %s:%s", name, username)
}

func (c *Client) CreatePermission(role, resource, permission string) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("action", permission)
	v.Add("resource", resource)
	v.Add("role", role)
	v.Add("accessToken", token)
	resp, err := http.PostForm(c.URL+api[c.APIVersion]["perm"], v)
	return checkErr(resp, err)
}

func (c *Client) DeletePermission(role, resource, permission string) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("action", permission)
	v.Add("resource", resource)
	v.Add("role", role)
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s%s?%s", c.URL, api[c.APIVersion]["perm"], v.Encode())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	return checkErr(resp, err)
}

func (c *Client) ListPermission() (*PermissionList, error) {
	if c.APIVersion == "v1" {
		return listResource[PermissionList](c, api[c.APIVersion]["list_perm"])
	}
	return listResource[PermissionListV3](c, api[c.APIVersion]["list_perm"])
}

func (c *Client) GetPermission(role, resource, action string) (*Permission, error) {
	perms, err := c.ListPermission()
	if err != nil {
		return nil, err
	}
	p := Permission{Role: role, Resource: resource, Action: action}
	if perms.Contains(p) {
		return &p, nil
	}
	return nil, fmt.Errorf("404 Not Found %s:%s:%s", role, resource, action)
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			return fmt.Errorf("%s %s %w", resp.Status, resp.Request.URL, err)
		}
		// no data or html data
		if len(data) == 0 || data[0] == '<' {
			return fmt.Errorf("%s %s", resp.Status, resp.Request.URL)
		}
		return fmt.Errorf("%s %s %s", resp.Status, resp.Request.URL, data)
	}
	return nil
}

func checkErr(resp *http.Response, httpErr error) error {
	if httpErr != nil {
		return httpErr
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

func decode(resp *http.Response, httpErr error, v any) error {
	if httpErr != nil {
		return httpErr
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return err
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

// type NacosErr struct {
// 	StatusCode int
// 	Err        error
// 	URL        string
// }

// func (e *NacosErr) Error() string {
// 	return fmt.Sprintf("%d %s: %s", e.StatusCode, e.URL, e.Err.Error())
// }

// func (e *NacosErr) Unwrap() error { return e.Err }
