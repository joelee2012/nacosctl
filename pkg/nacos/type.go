package nacos

type ConfigurationList struct {
	TotalCount     int              `json:"totalCount,omitempty"`
	PageNumber     int              `json:"pageNumber,omitempty"`
	PagesAvailable int              `json:"pagesAvailable,omitempty"`
	Items          []*Configuration `json:"pageItems"`
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
