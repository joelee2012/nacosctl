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
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

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

type ListTypes interface {
	User | Role | Permission | Configuration
}

type List[T ListTypes] struct {
	TotalCount     int  `json:"totalCount,omitempty"`
	PageNumber     int  `json:"pageNumber,omitempty"`
	PagesAvailable int  `json:"pagesAvailable,omitempty"`
	Items          []*T `json:"pageItems"`
}

func (o *List[T]) Contains(other T) bool {
	for _, it := range o.Items {
		if *it == other {
			return true
		}
	}
	return false
}

type ConfigurationList = List[Configuration]
type PermissionList = List[Permission]
type RoleList = List[Role]
type UserList = List[User]

func listResource[T User | Role | Permission](c *Client, endpoint string) (*List[T], error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	all := new(List[T])
	v := url.Values{}
	v.Add("search", "accurate")
	v.Add("accessToken", token)
	v.Add("pageNo", "1")
	v.Add("pageSize", "100")
	for {
		l := new(List[T])
		url := fmt.Sprintf("%s/%s?%s", c.URL, endpoint, v.Encode())
		resp, err := http.Get(url)
		if err := decode(resp, err, l); err != nil {
			return nil, err
		}
		all.Items = append(all.Items, l.Items...)
		if l.PagesAvailable == 0 || l.PagesAvailable == l.PageNumber {
			break
		}
		v.Set("pageNo", strconv.Itoa(l.PageNumber+1))
	}
	return all, nil
}

type V3List[T ListTypes] struct {
	Data *List[T] `json:"data"`
}

type ConfigurationListV3 = V3List[Configuration]
type PermissionListV3 = V3List[Permission]
type RoleListV3 = V3List[Role]
type UserListV3 = V3List[User]
