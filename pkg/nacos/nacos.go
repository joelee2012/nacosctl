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

func NewClient(url, user, password string) *Client {
	return &Client{
		URL:      url,
		User:     user,
		Password: password,
	}
}

func (c *Client) GetVersion() (string, error) {
	if c.State != nil {
		return c.Version, nil
	}
	resp, err := http.Get(c.URL + "/v1/console/server/state")
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
	resp, err := http.PostForm(c.URL+"/v1/auth/login", v)
	err = decode(resp, err, &c.Token)
	if err != nil {
		return "", err
	}
	c.ExpiredAt = now + c.TokenTTL
	return c.AccessToken, err
}

func (c *Client) ListNamespace() (*NsList, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s/v1/console/namespaces?%s", c.URL, v.Encode())
	resp, err := http.Get(url)
	namespaces := new(NsList)
	err = decode(resp, err, namespaces)
	return namespaces, err
}

type CreateNSOpts struct {
	Name        string
	Description string
	ID          string
}

func (c *Client) CreateNamespace(opts *CreateNSOpts) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("customNamespaceId", opts.ID)
	v.Add("namespaceName", opts.Name)
	v.Add("namespaceDesc", opts.Description)
	v.Add("accessToken", token)
	resp, err := http.PostForm(c.URL+"/v1/console/namespaces", v)
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
	url := fmt.Sprintf("%s/v1/console/namespaces?%s", c.URL, v.Encode())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	return checkErr(resp, err)
}

func (c *Client) UpdateNamespace(opts *CreateNSOpts) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("namespace", opts.ID)
	v.Add("namespaceShowName", opts.Name)
	v.Add("namespaceDesc", opts.Description)
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s/v1/console/namespaces?%s", c.URL, v.Encode())
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	return checkErr(resp, err)
}

func (c *Client) CreateOrUpdateNamespace(opts *CreateNSOpts) error {
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
			return ns, nil
		}
	}
	return nil, fmt.Errorf("404 Not Found %s", id)
}

type GetCSOpts struct {
	DataID      string
	Group       string
	NamespaceID string
}

func (c *Client) GetConfig(opts *GetCSOpts) (*Config, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Add("dataId", opts.DataID)
	v.Add("group", opts.Group)
	v.Add("namespaceId", opts.NamespaceID)
	v.Add("tenant", opts.NamespaceID)
	v.Add("show", "all")
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s/v1/cs/configs?%s", c.URL, v.Encode())
	resp, err := http.Get(url)
	config := new(Config)
	err = decode(resp, err, config)
	// if config not found, nacos server return 200 and empty response
	if err == io.EOF {
		return nil, fmt.Errorf("404 Not Found %s %w", url, err)
	}
	return config, err
}

type ListCSOpts struct {
	DataID      string
	Group       string
	Content     string
	AppName     string
	NamespaceID string
	PageNumber  int
	Tags        string
	PageSize    int
}

func (c *Client) ListConfig(opts *ListCSOpts) (*ConfigList, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Add("dataId", opts.DataID)
	v.Add("group", opts.Group)
	v.Add("appName", opts.AppName)
	v.Add("config_tags", opts.Tags)
	if opts.PageNumber == 0 {
		opts.PageNumber = 1
	}
	if opts.PageSize == 0 {
		opts.PageSize = 10
	}
	v.Add("pageNo", strconv.Itoa(opts.PageNumber))
	v.Add("pageSize", strconv.Itoa(opts.PageSize))
	v.Add("tenant", opts.NamespaceID)
	v.Add("search", "accurate")
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s/v1/cs/configs?%s", c.URL, v.Encode())
	configs := new(ConfigList)
	resp, err := http.Get(url)
	err = decode(resp, err, configs)
	return configs, err
}

func (c *Client) ListConfigInNs(namespace, group string) (*ConfigList, error) {
	nsCs := new(ConfigList)
	listOpts := ListCSOpts{PageNumber: 1, PageSize: 100, Group: group, NamespaceID: namespace}
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

func (c *Client) ListAllConfig() (*ConfigList, error) {
	allCs := new(ConfigList)
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

type CreateCSOpts struct {
	DataID      string
	Group       string
	Content     string
	AppName     string
	NamespaceID string
	Tags        string
	Type        string
	Desc        string
}

func (c *Client) CreateConfig(opts *CreateCSOpts) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("dataId", opts.DataID)
	v.Add("group", opts.Group)
	v.Add("content", opts.Content)
	v.Add("type", opts.Type)
	v.Add("tenant", opts.NamespaceID)
	v.Add("namespaceId", opts.NamespaceID)
	v.Add("appName", opts.AppName)
	v.Add("desc", opts.Desc)
	v.Add("config_tags", opts.Tags)
	v.Add("accessToken", token)
	resp, err := http.PostForm(c.URL+"/v1/cs/configs", v)
	return checkErr(resp, err)
}

type DeleteCSOpts = GetCSOpts

func (c *Client) DeleteConfig(opts *DeleteCSOpts) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("dataId", opts.DataID)
	v.Add("group", opts.Group)
	v.Add("tenant", opts.NamespaceID)
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s/v1/cs/configs?%s", c.URL, v.Encode())
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
	resp, err := http.PostForm(c.URL+"/v1/auth/users", v)
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
	url := fmt.Sprintf("%s/v1/auth/users?%s", c.URL, v.Encode())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	return checkErr(resp, err)
}

func (c *Client) ListUser() (*UserList, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	allUsers := new(UserList)
	v := url.Values{}
	v.Add("serach", "accurate")
	v.Add("accessToken", token)
	v.Add("pageNo", "1")
	v.Add("pageSize", "100")
	for {
		users := new(UserList)
		url := fmt.Sprintf("%s/v1/auth/users?%s", c.URL, v.Encode())
		resp, err := http.Get(url)
		if err := decode(resp, err, users); err != nil {
			return nil, err
		}
		allUsers.Items = append(allUsers.Items, users.Items...)
		if users.PagesAvailable == 0 || users.PagesAvailable == users.PageNumber {
			break
		}
		v.Set("pageNo", strconv.Itoa(users.PageNumber+1))
	}
	return allUsers, nil
}

func (c *Client) GetUser(name string) (*User, error) {
	users, err := c.ListUser()
	if err != nil {
		return nil, err
	}

	for _, user := range users.Items {
		if user.Name == name {
			return user, nil
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
	resp, err := http.PostForm(c.URL+"/v1/auth/roles", v)
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
	url := fmt.Sprintf("%s/v1/auth/roles?%s", c.URL, v.Encode())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	return checkErr(resp, err)
}

func (c *Client) ListRole() (*RoleList, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	allRoles := new(RoleList)
	v := url.Values{}
	v.Add("serach", "accurate")
	v.Add("accessToken", token)
	v.Add("pageNo", "1")
	v.Add("pageSize", "100")
	for {
		users := new(RoleList)
		url := fmt.Sprintf("%s/v1/auth/roles?%s", c.URL, v.Encode())
		resp, err := http.Get(url)
		if err := decode(resp, err, users); err != nil {
			return nil, err
		}
		allRoles.Items = append(allRoles.Items, users.Items...)
		if users.PagesAvailable == 0 || users.PagesAvailable == users.PageNumber {
			break
		}
		v.Set("pageNo", strconv.Itoa(users.PageNumber+1))
	}
	return allRoles, nil
}

func (c *Client) GetRole(name string) (*Role, error) {
	users, err := c.ListRole()
	if err != nil {
		return nil, err
	}

	for _, role := range users.Items {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, fmt.Errorf("404 Not Found %s", name)
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
	resp, err := http.PostForm(c.URL+"/v1/auth/permissions", v)
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
	url := fmt.Sprintf("%s/v1/auth/permissions?%s", c.URL, v.Encode())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	return checkErr(resp, err)
}

func (c *Client) ListPermission() (*PermissionList, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	allPerms := new(PermissionList)
	v := url.Values{}
	v.Add("serach", "accurate")
	v.Add("accessToken", token)
	v.Add("pageNo", "1")
	v.Add("pageSize", "100")
	for {
		perms := new(PermissionList)
		url := fmt.Sprintf("%s/v1/auth/permissions?%s", c.URL, v.Encode())
		resp, err := http.Get(url)
		if err := decode(resp, err, perms); err != nil {
			return nil, err
		}
		allPerms.Items = append(allPerms.Items, perms.Items...)
		if perms.PagesAvailable == 0 || perms.PagesAvailable == perms.PageNumber {
			break
		}
		v.Set("pageNo", strconv.Itoa(perms.PageNumber+1))
	}
	return allPerms, nil
}

func (c *Client) GetPermission(role, resource, permission string) (*Permission, error) {
	perms, err := c.ListPermission()
	if err != nil {
		return nil, err
	}

	for _, p := range perms.Items {
		if p.Role == role && p.Resource == resource && p.Permission == permission {
			return p, nil
		}
	}
	return nil, fmt.Errorf("404 Not Found %s:%s:%s", role, resource, permission)
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
