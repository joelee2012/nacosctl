package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type NsList struct {
	Code    int64        `json:"code"`
	Message interface{}  `json:"message"`
	Items   []*Namespace `json:"data"`
}

type Namespace struct {
	Name        string `json:"namespace"`
	ShowName    string `json:"namespaceShowName"`
	Desc        string `json:"namespaceDesc"`
	Quota       int64  `json:"quota"`
	ConfigCount int64  `json:"configCount"`
	Type        int64  `json:"type"`
}

type ConfigList struct {
	TotalCount     int64     `json:"totalCount"`
	PageNumber     int64     `json:"pageNumber"`
	PagesAvailable int64     `json:"pagesAvailable"`
	PageItems      []*Config `json:"pageItems"`
}

type Config struct {
	ID               string `json:"id"`
	DataID           string `json:"dataId"`
	Group            string `json:"group"`
	Content          string `json:"content"`
	Md5              string `json:"md5"`
	EncryptedDataKey string `json:"encryptedDataKey"`
	Tenant           string `json:"tenant"`
	AppName          string `json:"appName"`
	Type             string `json:"type"`
}

type Nacos struct {
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
}
type State struct {
	Version        string `json:"version"`
	StandaloneMode string `json:"standalone_mode"`
	FunctionMode   string `json:"function_mode"`
}

func NewNacos(url, user, password string) *Nacos {
	return &Nacos{
		URL:      url,
		User:     user,
		Password: password,
	}
}

func (n *Nacos) GetVersion() (string, error) {
	if n.State != nil {
		return n.Version, nil
	}
	resp, err := http.Get(n.URL + "/nacos/v1/console/server/state")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&n.State); err != nil {
		return "", err
	}
	return n.Version, nil
}

func (n *Nacos) GetToken() (string, error) {
	if n.Token != nil {
		return n.AccessToken, nil
	}
	v := url.Values{}
	v.Add("username", n.User)
	v.Add("password", n.Password)
	resp, err := http.PostForm(n.URL+"/nacos/v1/auth/login", v)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&n.Token); err != nil {
		return "", err
	}
	return n.AccessToken, nil
}

func (n *Nacos) ListNamespace() (*NsList, error) {
	token, err := n.GetToken()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Add("accessToken", token)
	url := fmt.Sprintf("%s/nacos/v1/console/namespaces?%s", n.URL, v.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	dec := json.NewDecoder(resp.Body)
	namespaces := &NsList{}
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

func (n *Nacos) CreateNamespace(opts *CreateNSOpts) error {
	token, err := n.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("customNamespaceId", opts.ID)
	v.Add("namespaceName", opts.Name)
	v.Add("namespaceDesc", opts.Desc)
	v.Add("accessToken", token)
	v.Add("username", n.User)
	resp, err := http.PostForm(n.URL+"/nacos/v1/console/namespaces", v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create namespace/%s with error %s", opts.Name, resp.Status)
	}
	return nil
}

func (n *Nacos) DeleteNamespace(id string) error {
	token, err := n.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("namespaceId", id)
	v.Add("accessToken", token)
	v.Add("username", n.User)
	url := fmt.Sprintf("%s/nacos/v1/console/namespaces?%s", n.URL, v.Encode())
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

func (n *Nacos) UpdateNamespace(opts *CreateNSOpts) error {
	token, err := n.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("namespace", opts.ID)
	v.Add("namespaceShowName", opts.Name)
	v.Add("namespaceDesc", opts.Desc)
	v.Add("accessToken", token)
	v.Add("username", n.User)

	url := fmt.Sprintf("%s/nacos/v1/console/namespaces?%s", n.URL, v.Encode())
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

func (n *Nacos) ListConfig(opts *ListCSOpts) (*ConfigList, error) {
	token, err := n.GetToken()
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
	v.Add("search", "accurate")
	v.Add("accessToken", token)
	v.Add("username", n.User)
	url := fmt.Sprintf("%s/nacos/v1/cs/configs?%s", n.URL, v.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var configs ConfigList
	if err := dec.Decode(&configs); err != nil {
		return nil, err
	}

	return &configs, nil
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

func (n *Nacos) CreateConfig(opts *CreateCSOpts) error {
	token, err := n.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("dataId", opts.DataID)
	v.Add("group", opts.Group)
	v.Add("content", opts.Content)
	v.Add("type", opts.Type)
	v.Add("tenant", opts.Tenant)
	v.Add("accessToken", token)
	v.Add("username", n.User)
	resp, err := http.PostForm(n.URL+"/nacos/v1/cs/configs", v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create configuration/%s with %s", opts.DataID, resp.Status)
	}
	return nil
}
func (n *Nacos) DeleteConfig(opts *CreateCSOpts) error {
	token, err := n.GetToken()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("dataId", opts.DataID)
	v.Add("group", opts.Group)
	v.Add("tenant", opts.Tenant)
	v.Add("accessToken", token)
	v.Add("username", n.User)
	url := fmt.Sprintf("%s/nacos/v1/cs/configs?%s", n.URL, v.Encode())
	fmt.Println(url)
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
