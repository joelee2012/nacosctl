package nacos

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
var namespace = `
{
	"namespace": "test",
	"namespaceShowName": "Test",
	"namespaceDesc": "Test namespace",
	"quota": 100,
	"configCount": 10,
	"type": 0
}`

var role = `{"role": "ROLE_ADMIN", "username": "nacos"}`

var user = `{"username": "user1", "password": "$2a$10$C3B9EQgp93M6mvXwXiCebe1T9HvxGRj29x2dHIYCH.bUCdbJcrugO"}`
var permission = `{"role": "ROLE_ADMIN", "resource": "backend:*:*", "action": "rw"}`

func newV1Data(s string) string {
	return fmt.Sprintf(`{
  "totalCount": 1,
  "pageNumber": 1,
  "pagesAvailable": 0,
  "pageItems": [
    %s
  ]
}`, s)
}

func newV3Data(s string) string {
	return fmt.Sprintf(`{
  "code": 0,
  "message": "success",
  "data": %s
}
`, s)
}

var csList = newV1Data(config)

var csListV3 = newV3Data(csList)

var configV3 = newV3Data(config)

var nsList = fmt.Sprintf(`
{
  "code": 200,
  "message": "success",
  "data": [
    %s
  ]
}
`, namespace)

var userList = newV1Data(user)
var userListV3 = newV3Data(userList)

var roleList = newV1Data(role)
var roleListV3 = newV3Data(roleList)

var permList = newV1Data(permission)
var permListV3 = newV3Data(permList)

func TestNewClient(t *testing.T) {
	c := NewClient("http://localhost:8848", "user", "password")
	assert.Equal(t, "http://localhost:8848", c.URL)
	assert.Equal(t, "user", c.User)
	assert.Equal(t, "password", c.Password)
}

func startServer() (*httptest.Server, *Client) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		switch r.URL.Path {
		case "/v1/console/namespaces", "/v3/console/core/namespace/list":
			if r.URL.Query().Get("show") == "all" {
				w.Write([]byte(namespace))
			} else {
				w.Write([]byte(nsList))
			}
		case "/v1/cs/configs":
			if r.URL.Query().Get("show") == "all" {
				w.Write([]byte(config))
			} else {
				w.Write([]byte(csList))
			}
		case "/v3/console/cs/config":
			w.Write([]byte(configV3))
		case "/v3/console/cs/config/list":
			w.Write([]byte(csListV3))
		case "/v1/console/server/state":
			w.Write([]byte(`{"version": "1.0.0"}`))
		case "/v3/console/server/state":
			w.Write([]byte(`{"version": "3.0.0"}`))
		case "/v1/auth/login", "/v3/auth/user/login":
			w.Write([]byte(`{"accessToken": "test-token", "tokenTtl": 3600, "globalAdmin": true}`))
		case "/v1/auth/users":
			w.Write([]byte(userList))
		case "/v3/auth/user/list":
			w.Write([]byte(userListV3))
		case "/v1/auth/roles":
			w.Write([]byte(roleList))
		case "/v3/auth/role/list":
			w.Write([]byte(roleListV3))
		case "/v1/auth/permissions":
			w.Write([]byte(permList))
		case "/v3/auth/permission/list":
			w.Write([]byte(permListV3))
		}
	}))
	c := NewClient(ts.URL, "user", "password")
	return ts, c
}

var apiTests = []struct {
	apiVersion  string
	expectValue string
}{
	{apiVersion: "v1", expectValue: "1.0.0"},
	{apiVersion: "v3", expectValue: "3.0.0"},
}

func TestGetVersion(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()
	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			c.State = nil
			version, err := c.GetVersion()
			if assert.NoError(t, err) {
				assert.Equal(t, tt.expectValue, version)
			}
		})

	}
}

func startAccClient() *Client {
	return NewClient(os.Getenv("NACOS_HOST"), os.Getenv("NACOS_USER"), os.Getenv("NACOS_PASSWORD"))
}

func TestDetectAPIVersion(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()
	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			// mock api version list
			apiVersions = []string{tt.apiVersion}
			c.DetectAPIVersion()
			assert.Equal(t, tt.apiVersion, c.APIVersion)
		})

	}
}

func TestAccDetectAPIVersion(t *testing.T) {
	if os.Getenv("ACC") != "true" {
		t.Skip("skip as ACC != true ")
	}
	ts := startAccClient()
	ts.DetectAPIVersion()
	assert.Equal(t, "v1", ts.APIVersion)
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

}

func TestListNamespace(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			ns, err := c.ListNamespace()
			if assert.NoError(t, err) {
				assert.Equal(t, 1, len(ns.Items))
				assert.Equal(t, "test", ns.Items[0].ID)
			}
		})
	}
}

func TestCreateNamespace(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			err := c.CreateNamespace(&CreateNsOpts{Name: "test", Description: "Test namespace", ID: "test-id"})
			assert.NoError(t, err)
		})
	}
}

func TestGetNamespace(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()
	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			n, err := c.GetNamespace("test")
			if assert.NoError(t, err) {
				assert.Equal(t, "test", n.ID)
			}
		})
	}
}

func TestDeleteNamespace(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()
	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			err := c.DeleteNamespace("test-id")
			assert.NoError(t, err)
		})
	}
}

func TestUpdateNamespace(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()
	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			err := c.UpdateNamespace(&CreateNsOpts{Name: "test", Description: "Test namespace", ID: "test-id"})
			assert.NoError(t, err)
		})
	}
}
func TestCreateOrUpdateNamespace(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	tests := []struct {
		name string
		data CreateNsOpts
	}{
		{name: "Create", data: CreateNsOpts{Name: "test", Description: "Test namespace", ID: "test"}},
		{name: "Update", data: CreateNsOpts{Name: "test-id", Description: "Test namespace", ID: "test-id"}},
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
	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			cfg, err := c.GetConfig(&GetCfgOpts{DataID: "test", Group: "DEFAULT_GROUP"})
			if assert.NoError(t, err) {
				assert.Equal(t, "test", cfg.DataID)
			}
		})
	}
}

func TestListConfig(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()
	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			cfgs, err := c.ListConfig(&ListCfgOpts{DataID: "test", Group: "DEFAULT_GROUP", PageNumber: 1, PageSize: 10})
			if assert.NoError(t, err) {
				assert.Equal(t, 1, len(cfgs.Items))
				assert.Equal(t, "test", cfgs.Items[0].DataID)
			}
		})
	}
}

func TestListConfigInNs(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			cfgs, err := c.ListConfigInNs("test", "DEFAULT_GROUP")
			if assert.NoError(t, err) {
				assert.Equal(t, 1, len(cfgs.Items))
				assert.Equal(t, "test", cfgs.Items[0].DataID)
			}
		})
	}
}

func TestListAllConfig(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()
	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			cfgs, err := c.ListAllConfig()
			if assert.NoError(t, err) {
				assert.Equal(t, 1, len(cfgs.Items))
				assert.Equal(t, "test", cfgs.Items[0].DataID)
			}
		})
	}
}

func TestCreateConfig(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.CreateConfig(&CreateCfgOpts{DataID: "test", Group: "DEFAULT_GROUP", Content: "test content", NamespaceID: "test-tenant", Type: "properties"})
	assert.NoError(t, err)
}

func TestDeleteConfig(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()

	err := c.DeleteConfig(&DeleteCfgOpts{DataID: "test", Group: "DEFAULT_GROUP", NamespaceID: "test-tenant"})
	assert.NoError(t, err)
}

func TestListUser(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()
	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			users, err := c.ListUser()
			if assert.NoError(t, err) {
				assert.Equal(t, "user1", users.Items[0].Name)
				// assert.Equal(t, "user2", users.Items[1].Name)
			}
		})
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
	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			roles, err := c.ListRole()
			if assert.NoError(t, err) {
				assert.Equal(t, "ROLE_ADMIN", roles.Items[0].Name)
				assert.Equal(t, "nacos", roles.Items[0].Username)
			}
		})
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

	role, err := c.GetRole("ROLE_ADMIN", "nacos")
	if assert.NoError(t, err) {
		assert.Equal(t, "ROLE_ADMIN", role.Name)
	}
}

func TestListPermission(t *testing.T) {
	ts, c := startServer()
	defer ts.Close()
	for _, tt := range apiTests {
		t.Run(tt.apiVersion, func(t *testing.T) {
			c.APIVersion = tt.apiVersion
			perms, err := c.ListPermission()
			if assert.NoError(t, err) {
				assert.Equal(t, "ROLE_ADMIN", perms.Items[0].Role)
				assert.Equal(t, "backend:*:*", perms.Items[0].Resource)
				assert.Equal(t, "rw", perms.Items[0].Action)
			}
		})
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
