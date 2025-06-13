package nacos

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type NsList struct {
	// Code    int          `json:"code,omitempty"`
	// Message interface{}  `json:"message,omitempty"`
	Items []*Namespace `json:"data"`
}

type Namespace struct {
	Name        string `json:"namespace"`
	ShowName    string `json:"namespaceShowName"`
	Desc        string `json:"namespaceDesc"`
	Quota       int    `json:"quota"`
	ConfigCount int    `json:"configCount"`
	Type        int    `json:"type"`
}

type ConfigList struct {
	TotalCount     int       `json:"totalCount,omitempty"`
	PageNumber     int       `json:"pageNumber,omitempty"`
	PagesAvailable int       `json:"pagesAvailable,omitempty"`
	Items          []*Config `json:"pageItems"`
}

type Config struct {
	ID               string `json:"id"`
	DataID           string `json:"dataId"`
	Group            string `json:"group"`
	Content          string `json:"content"`
	Tenant           string `json:"tenant"`
	Type             string `json:"type"`
	Md5              string `json:"md5,omitempty"`
	EncryptedDataKey string `json:"encryptedDataKey,omitempty"`
	AppName          string `json:"appName,omitempty"`
	CreateTime       int64  `json:"createTime,omitempty"`
	ModifyTime       int64  `json:"modifyTime,omitempty"`
	Desc             string `json:"desc,omitempty"`
	Tags             string `json:"configTags,omitempty"`
}

type Client struct {
	URL        string
	User       string
	Password   string
	APIVersion string
	*Token
	*State
	// client *http.Client
}
type Token struct {
	AccessToken string `json:"accessToken"`
	TokenTTL    int    `json:"tokenTtl"`
	GlobalAdmin bool   `json:"globalAdmin"`
	Username    string `json:"username"`
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

// func (n *Nacos) Client() *http.Client {
// 	if n.client == nil {
// 		n.client = &http.Client{}
// 	}
// 	return n.client
// }

func (c *Client) GetVersion() (string, error) {
	if c.State != nil {
		return c.Version, nil
	}
	resp, err := http.Get(c.URL + "/nacos/v1/console/server/state")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&c.State); err != nil {
		return "", err
	}
	return c.Version, nil
}

func (c *Client) GetToken() (string, error) {
	if c.Token != nil {
		return c.AccessToken, nil
	}
	v := url.Values{}
	v.Add("username", c.User)
	v.Add("password", c.Password)
	resp, err := http.PostForm(c.URL+"/nacos/v1/auth/login", v)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&c.Token); err != nil {
		return "", err
	}
	return c.AccessToken, nil
}

func (c *Client) ListNamespace() (*NsList, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s/nacos/v1/console/namespaces?%s", c.URL, v.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	dec := json.NewDecoder(resp.Body)
	namespaces := new(NsList)
	if err := dec.Decode(namespaces); err != nil {
		return nil, err
	}
	return namespaces, nil
}

type CreateNSOpts struct {
	Name string
	Desc string
	ID   string
}

func (c *Client) CreateNamespace(opts *CreateNSOpts) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("customNamespaceId", opts.ID)
	v.Add("namespaceName", opts.Name)
	v.Add("namespaceDesc", opts.Desc)
	v.Add("accessToken", token)
	v.Add("username", c.User)
	resp, err := http.PostForm(c.URL+"/nacos/v1/console/namespaces", v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create namespace/%s with error %s", opts.Name, resp.Status)
	}
	return nil
}

func (c *Client) DeleteNamespace(id string) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("namespaceId", id)
	v.Add("accessToken", token)
	v.Add("username", c.User)
	url := fmt.Sprintf("%s/nacos/v1/console/namespaces?%s", c.URL, v.Encode())
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete namespace/%s with error %s", id, resp.Status)
	}
	return nil
}

func (c *Client) UpdateNamespace(opts *CreateNSOpts) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("namespace", opts.ID)
	v.Add("namespaceShowName", opts.Name)
	v.Add("namespaceDesc", opts.Desc)
	v.Add("accessToken", token)
	v.Add("username", c.User)

	url := fmt.Sprintf("%s/nacos/v1/console/namespaces?%s", c.URL, v.Encode())
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update namespace/%s status code: %d", opts.ID, resp.StatusCode)
	}
	return nil
}

func (c *Client) CreateOrUpdateNamespace(opts *CreateNSOpts) error {
	nsList, err := c.ListNamespace()
	if err != nil {
		return err
	}
	for _, ns := range nsList.Items {
		if ns.Name == opts.ID {
			return c.UpdateNamespace(opts)
		}
	}
	return c.CreateNamespace(opts)
}

type ListCSOpts struct {
	DataId     string
	Group      string
	Content    string
	AppName    string
	Tenant     string
	PageNumber int
	Tags       string
	PageSize   int
}

func (c *Client) ListConfig(opts *ListCSOpts) (*ConfigList, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Add("dataId", opts.DataId)
	v.Add("group", opts.Group)
	v.Add("appName", opts.AppName)
	v.Add("config_tags", opts.Tags)
	v.Add("pageNo", strconv.Itoa(opts.PageNumber))
	v.Add("pageSize", strconv.Itoa(opts.PageSize))
	v.Add("tenant", opts.Tenant)
	// v.Add("show", "all")
	v.Add("search", "accurate")
	v.Add("accessToken", token)
	v.Add("username", c.User)
	url := fmt.Sprintf("%s/nacos/v1/cs/configs?%s", c.URL, v.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	configs := new(ConfigList)
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func (c *Client) ListConfigInNs(namespace, group string) (*ConfigList, error) {
	nsCs := new(ConfigList)
	listOpts := ListCSOpts{PageNumber: 1, PageSize: 100, Group: group, Tenant: namespace}
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
		cs, err := c.ListConfigInNs(ns.Name, "")
		if err != nil {
			return nil, err
		}
		allCs.Items = append(allCs.Items, cs.Items...)
	}
	return allCs, nil
}

type CreateCSOpts struct {
	DataID  string
	Group   string
	Content string
	AppName string
	Tenant  string
	Tags    string
	Type    string
	Desc    string
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
	v.Add("tenant", opts.Tenant)
	v.Add("namespaceId", opts.Tenant)
	v.Add("appName", opts.AppName)
	v.Add("desc", opts.Desc)
	v.Add("config_tags", opts.Tags)
	v.Add("accessToken", token)
	v.Add("username", c.User)
	resp, err := http.PostForm(c.URL+"/nacos/v1/cs/configs", v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create configuration/%s with %s", opts.DataID, resp.Status)
	}
	return nil
}

type DeleteCSOpts struct {
	DataID string
	Group  string
	Tenant string
}

func (c *Client) DeleteConfig(opts *DeleteCSOpts) error {
	token, err := c.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("dataId", opts.DataID)
	v.Add("group", opts.Group)
	v.Add("tenant", opts.Tenant)
	v.Add("accessToken", token)
	v.Add("username", c.User)
	url := fmt.Sprintf("%s/nacos/v1/cs/configs?%s", c.URL, v.Encode())
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete configuration/%s with %s", opts.DataID, resp.Status)
	}
	return nil
}
