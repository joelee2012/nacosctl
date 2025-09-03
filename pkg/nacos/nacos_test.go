package nacos

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var csList = `
{
  "totalCount": 1,
  "pageNumber": 1,
  "pagesAvailable": 0,
  "pageItems": [
    {
      "id": "1",
      "dataId": "test",
      "group": "DEFAULT_GROUP",
      "content": "test content",
      "md5": "test-md5",
      "encryptedDataKey": "test-key",
      "tenant": "test-tenant",
      "appName": "test-app",
      "type": "properties"
    }
  ]
}
`

var config = `
{
	"id": "1",
	"dataId": "test",
	"group": "DEFAULT_GROUP",
	"content": "test content",
	"md5": "test-md5",
	"encryptedDataKey": "test-key",
	"tenant": "test-tenant",
	"appName": "test-app",
	"type": "properties"
}
`
var nsList = `
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "namespace": "test",
      "namespaceShowName": "Test",
      "namespaceDesc": "Test namespace",
      "quota": 100,
      "configCount": 10,
      "type": 0
    }
  ]
}
`
var namespace = `
{
	"namespace": "test",
	"namespaceShowName": "Test",
	"namespaceDesc": "Test namespace",
	"quota": 100,
	"configCount": 10,
	"type": 0
}
`
var users = `
{
  "totalCount": 2,
  "pageNumber": 1,
  "pagesAvailable": 1,
  "pageItems": [
    {
      "username": "user1",
      "password": "$2a$10$C3B9EQgp93M6mvXwXiCebe1T9HvxGRj29x2dHIYCH.bUCdbJcrugO"
    },
    {
      "username": "user2",
      "password": "$2a$10$OHWIUaiy9cC8wHxANJ7j/O6CDNe1fy5WUD/6vUA2TdeWNjZLPUA.C"
    }
  ]
}
`

var roles = `
{
  "totalCount": 1,
  "pageNumber": 1,
  "pagesAvailable": 1,
  "pageItems": [
    {
      "role": "ROLE_ADMIN",
      "username": "nacos"
    }
  ]
}
`

var permissions = `
{
  "totalCount": 1,
  "pageNumber": 1,
  "pagesAvailable": 1,
  "pageItems": [
    {
      "role": "ROLE_ADMIN",
      "resource": "backend:*:*",
      "action": "rw"
    }
  ]
}
`

func TestNewClient(t *testing.T) {
	c := NewClient("http://localhost:8848", "user", "password")
	assert.Equal(t, "http://localhost:8848", c.URL)
	assert.Equal(t, "user", c.User)
	assert.Equal(t, "password", c.Password)
}

func startServer() (*httptest.Server, *Client) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.URL.Path == "/v1/console/namespaces" {
			if r.URL.Query().Get("show") == "all" {
				w.Write([]byte(namespace))
			} else {
				w.Write([]byte(nsList))
			}
		} else if r.URL.Path == "/v1/cs/configs" {
			if r.URL.Query().Get("show") == "all" {
				w.Write([]byte(config))
			} else {
				w.Write([]byte(csList))
			}
		} else if r.URL.Path == "/v1/console/server/state" {
			w.Write([]byte(`{"version": "1.0.0"}`))
		} else if r.URL.Path == "/v1/auth/login" {
			w.Write([]byte(`{"accessToken": "test-token", "tokenTtl": 3600, "globalAdmin": true}`))
		} else if r.URL.Path == "/v1/auth/users" {
			w.Write([]byte(users))
		} else if r.URL.Path == "/v1/auth/roles" {
			w.Write([]byte(roles))
		} else if r.URL.Path == "/v1/auth/permissions" {
			w.Write([]byte(permissions))
		}
	}))
	c := NewClient(ts.URL, "user", "password")
	return ts, c
}
func TestGetVersion(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()
	version, err := c.GetVersion()
	if assert.NoError(t, err) {
		assert.Equal(t, "1.0.0", version)
	}
}

func TestGetToken(t *testing.T) {
	okServer, _ := startServer()
	defer okServer.Close()
	nokServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`<html><body>404 Not Found</body></html>`))
	}))
	defer nokServer.Close()
	tests := []struct {
		name    string
		server  *httptest.Server
		wantErr bool
	}{
		{name: "OK", server: okServer, wantErr: false},
		{name: "NOK", server: nokServer, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.server.URL, "user", "password")
			token, err := c.GetToken()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, "", token)
				assert.Equal(t, err.Error(), fmt.Sprintf("404 Not Found %s/v1/auth/login", tt.server.URL))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "test-token", token)
			}
		})
	}
	// c := NewClient(okServer.URL, "user", "password")
	// token, err := c.GetToken()
	// if assert.NoError(t, err) {
	// 	assert.Equal(t, "test-token", token)
	// }
	// c = NewClient("http://wrong.context:8080", "user", "password")
	// token, err = c.GetToken()
	// if assert.Error(t, err) {
	// 	assert.Equal(t, "", token)
	// }
}

func TestListNamespace(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	namespaces, err := c.ListNamespace()
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(namespaces.Items))
		assert.Equal(t, "test", namespaces.Items[0].ID)
	}
}

func TestCreateNamespace(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.CreateNamespace(&CreateNSOpts{Name: "test", Description: "Test namespace", ID: "test-id"})
	assert.NoError(t, err)
}

func TestGetNamespace(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	ns, err := c.GetNamespace("test")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", ns.ID)
	}
}

func TestDeleteNamespace(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.DeleteNamespace("test-id")
	assert.NoError(t, err)
}

func TestUpdateNamespace(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.UpdateNamespace(&CreateNSOpts{Name: "test", Description: "Test namespace", ID: "test-id"})
	assert.NoError(t, err)
}
func TestCreateOrUpdateNamespace(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	tests := []struct {
		name string
		data CreateNSOpts
	}{
		{name: "Create", data: CreateNSOpts{Name: "test", Description: "Test namespace", ID: "test"}},
		{name: "Update", data: CreateNSOpts{Name: "test-id", Description: "Test namespace", ID: "test-id"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.CreateOrUpdateNamespace(&tt.data)
			assert.NoError(t, err)
		})
	}
}

func TestGetConfig(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	config, err := c.GetConfig(&GetCSOpts{DataID: "test", Group: "DEFAULT_GROUP"})
	if assert.NoError(t, err) {
		assert.Equal(t, "test", config.DataID)
	}
}

func TestListConfig(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	configs, err := c.ListConfig(&ListCSOpts{DataID: "test", Group: "DEFAULT_GROUP", PageNumber: 1, PageSize: 10})
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(configs.Items))
		assert.Equal(t, "test", configs.Items[0].DataID)
	}
}

func TestListConfigInNs(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	configs, err := c.ListConfigInNs("test", "DEFAULT_GROUP")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(configs.Items))
		assert.Equal(t, "test", configs.Items[0].DataID)
	}
}

func TestListAllConfig(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	configs, err := c.ListAllConfig()
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(configs.Items))
		assert.Equal(t, "test", configs.Items[0].DataID)
	}
}

func TestCreateConfig(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.CreateConfig(&CreateCSOpts{DataID: "test", Group: "DEFAULT_GROUP", Content: "test content", NamespaceID: "test-tenant", Type: "properties"})
	assert.NoError(t, err)
}

func TestDeleteConfig(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.DeleteConfig(&DeleteCSOpts{DataID: "test", Group: "DEFAULT_GROUP", NamespaceID: "test-tenant"})
	assert.NoError(t, err)
}

func TestListUser(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	users, err := c.ListUser()
	if assert.NoError(t, err) {
		assert.Equal(t, "user1", users.Items[0].Name)
		assert.Equal(t, "user2", users.Items[1].Name)
	}
}

func TestCreateUser(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.CreateUser("user3", "password")
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.DeleteUser("user3")
	assert.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	user, err := c.GetUser("user1")
	if assert.NoError(t, err) {
		assert.Equal(t, "user1", user.Name)
	}
}

func TestListRole(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	users, err := c.ListRole()
	if assert.NoError(t, err) {
		assert.Equal(t, "ROLE_ADMIN", users.Items[0].Name)
		assert.Equal(t, "nacos", users.Items[0].Username)
	}
}

func TestCreateRole(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.CreateRole("role1", "user1")
	assert.NoError(t, err)
}

func TestDeleteRole(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.DeleteRole("role1", "user1")
	assert.NoError(t, err)
}

func TestGetRole(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	user, err := c.GetRole("ROLE_ADMIN", "nacos")
	if assert.NoError(t, err) {
		assert.Equal(t, "ROLE_ADMIN", user.Name)
	}
}

func TestListPermission(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	users, err := c.ListPermission()
	if assert.NoError(t, err) {
		assert.Equal(t, "ROLE_ADMIN", users.Items[0].Role)
		assert.Equal(t, "backend:*:*", users.Items[0].Resource)
		assert.Equal(t, "rw", users.Items[0].Action)
	}
}

func TestCreatePermission(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.CreatePermission("ROLE_ADMIN", "backend:*:*", "rw")
	assert.NoError(t, err)
}

func TestDeletePermission(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.DeletePermission("ROLE_ADMIN", "backend:*:*", "rw")
	assert.NoError(t, err)
}

func TestGetPermission(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	perm, err := c.GetPermission("ROLE_ADMIN", "backend:*:*", "rw")
	if assert.NoError(t, err) {
		assert.Equal(t, "ROLE_ADMIN", perm.Role)
		assert.Equal(t, "backend:*:*", perm.Resource)
		assert.Equal(t, "rw", perm.Action)
	}
}
