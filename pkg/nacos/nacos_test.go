package nacos

import (
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

func TestNewClient(t *testing.T) {
	c := NewClient("http://localhost:8848", "user", "password")
	assert.Equal(t, "http://localhost:8848", c.URL)
	assert.Equal(t, "user", c.User)
	assert.Equal(t, "password", c.Password)
}

func startServer() *httptest.Server {
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
		}
	}))
	return ts
}
func TestGetVersion(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	version, err := c.GetVersion()
	if assert.NoError(t, err) {
		assert.Equal(t, "1.0.0", version)
	}
}

func TestGetToken(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	token, err := c.GetToken()
	if assert.NoError(t, err) {
		assert.Equal(t, "test-token", token)
	}
}

func TestListNamespace(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	namespaces, err := c.ListNamespace()
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(namespaces.Items))
		assert.Equal(t, "test", namespaces.Items[0].ID)
	}
}

func TestCreateNamespace(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	n := NewClient(ts.URL, "user", "password")
	err := n.CreateNamespace(&CreateNSOpts{Name: "test", Description: "Test namespace", ID: "test-id"})
	assert.NoError(t, err)
}

func TestGetNamespace(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	n := NewClient(ts.URL, "user", "password")
	ns, err := n.GetNamespace("test")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", ns.ID)
	}
}

func TestDeleteNamespace(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	err := c.DeleteNamespace("test-id")
	assert.NoError(t, err)
}

func TestUpdateNamespace(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	err := c.UpdateNamespace(&CreateNSOpts{Name: "test", Description: "Test namespace", ID: "test-id"})
	assert.NoError(t, err)
}
func TestCreateOrUpdateNamespace(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
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
	ts := startServer()
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	config, err := c.GetConfig(&ListCSOpts{DataID: "test", Group: "DEFAULT_GROUP", PageNumber: 1, PageSize: 10})
	if assert.NoError(t, err) {
		assert.Equal(t, "test", config.DataID)
	}
}

func TestListConfig(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	configs, err := c.ListConfig(&ListCSOpts{DataID: "test", Group: "DEFAULT_GROUP", PageNumber: 1, PageSize: 10})
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(configs.Items))
		assert.Equal(t, "test", configs.Items[0].DataID)
	}
}

func TestListConfigInNs(t *testing.T) {
	ts := startServer()
	defer ts.Close()
	c := NewClient(ts.URL, "user", "password")
	configs, err := c.ListConfigInNs("test", "DEFAULT_GROUP")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(configs.Items))
		assert.Equal(t, "test", configs.Items[0].DataID)
	}
}

func TestListAllConfig(t *testing.T) {
	ts := startServer()
	defer ts.Close()
	c := NewClient(ts.URL, "user", "password")
	configs, err := c.ListAllConfig()
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(configs.Items))
		assert.Equal(t, "test", configs.Items[0].DataID)
	}
}

func TestCreateConfig(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	c.Token = &Token{AccessToken: "test-token"}
	err := c.CreateConfig(&CreateCSOpts{DataID: "test", Group: "DEFAULT_GROUP", Content: "test content", NamespaceID: "test-tenant", Type: "properties"})
	assert.NoError(t, err)
}

func TestDeleteConfig(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	c.Token = &Token{AccessToken: "test-token"}
	err := c.DeleteConfig(&DeleteCSOpts{DataID: "test", Group: "DEFAULT_GROUP", NamespaceID: "test-tenant"})
	assert.NoError(t, err)
}
