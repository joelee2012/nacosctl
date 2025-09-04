package nacos

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Configuration struct {
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

type NamespaceList struct {
	// Code    int          `json:"code,omitempty"`
	// Message interface{}  `json:"message,omitempty"`
	Items []*Namespace `json:"data"`
}

type Namespace struct {
	ID          string `json:"namespace"`
	Name        string `json:"namespaceShowName"`
	Description string `json:"namespaceDesc"`
	Quota       int    `json:"quota,omitempty"`
	ConfigCount int    `json:"configCount,omitempty"`
	Type        int    `json:"type,omitempty"`
}

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type Role struct {
	Name     string `json:"role"`
	Username string `json:"username"`
}

type Permission struct {
	Role     string `json:"role"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type ResourceTypes interface {
	User | Role | Permission | Configuration
}

type ObjectList[T ResourceTypes] struct {
	TotalCount     int  `json:"totalCount,omitempty"`
	PageNumber     int  `json:"pageNumber,omitempty"`
	PagesAvailable int  `json:"pagesAvailable,omitempty"`
	Items          []*T `json:"pageItems"`
}

type ConfigurationList = ObjectList[Configuration]
type UserList = ObjectList[User]
type RoleList = ObjectList[Role]
type PermissionList = ObjectList[Permission]

func listResource[T ResourceTypes](c *Client, endpoint string) (*ObjectList[T], error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	all := new(ObjectList[T])
	v := url.Values{}
	v.Add("search", "accurate")
	v.Add("accessToken", token)
	v.Add("pageNo", "1")
	v.Add("pageSize", "100")
	for {
		perms := new(ObjectList[T])
		url := fmt.Sprintf("%s/%s?%s", c.URL, endpoint, v.Encode())
		resp, err := http.Get(url)
		if err := decode(resp, err, perms); err != nil {
			return nil, err
		}
		all.Items = append(all.Items, perms.Items...)
		if perms.PagesAvailable == 0 || perms.PagesAvailable == perms.PageNumber {
			break
		}
		v.Set("pageNo", strconv.Itoa(perms.PageNumber+1))
	}
	return all, nil
}
