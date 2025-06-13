package nacos

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCsList = `
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
var testNsList = `
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

func TestNewClient(t *testing.T) {
	c := NewClient("http://localhost:8848", "user", "password")
	assert.Equal(t, "http://localhost:8848", c.URL)
	assert.Equal(t, "user", c.User)
	assert.Equal(t, "password", c.Password)
}

func TestGetVersion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version": "1.0.0"}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	version, err := c.GetVersion()
	if assert.NoError(t, err) {
		assert.Equal(t, "1.0.0", version)
	}
}

func TestGetToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"accessToken": "test-token", "tokenTtl": 3600, "globalAdmin": true}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	token, err := c.GetToken()
	if assert.NoError(t, err) {
		assert.Equal(t, "test-token", token)
	}
}

func TestListNamespace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testNsList))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	namespaces, err := c.ListNamespace()
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(namespaces.Items))
		assert.Equal(t, "test", namespaces.Items[0].Name)
	}
}

func TestCreateNamespace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewClient(ts.URL, "user", "password")
	n.Token = &Token{AccessToken: "test-token"}
	err := n.CreateNamespace(&CreateNSOpts{Name: "test", Desc: "Test namespace", ID: "test-id"})
	assert.NoError(t, err)
}

func TestDeleteNamespace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	c.Token = &Token{AccessToken: "test-token"}
	err := c.DeleteNamespace("test-id")
	assert.NoError(t, err)
}

func TestUpdateNamespace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	c.Token = &Token{AccessToken: "test-token"}
	err := c.UpdateNamespace(&CreateNSOpts{Name: "test", Desc: "Test namespace", ID: "test-id"})
	assert.NoError(t, err)
}
func TestCreateOrUpdateNamespace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(testNsList))
		}
		if r.Method == "POST" || r.Method == "PUT" {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	c.Token = &Token{AccessToken: "test-token"}
	tests := []struct {
		name string
		data CreateNSOpts
	}{
		{name: "Create", data: CreateNSOpts{Name: "test", Desc: "Test namespace", ID: "test"}},
		{name: "Update", data: CreateNSOpts{Name: "test-id", Desc: "Test namespace", ID: "test-id"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.CreateOrUpdateNamespace(&tt.data)
			assert.NoError(t, err)
		})
	}
}

func TestListConfig(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testCsList))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	configs, err := c.ListConfig(&ListCSOpts{DataId: "test", Group: "DEFAULT_GROUP", PageNumber: 1, PageSize: 10})
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(configs.Items))
		assert.Equal(t, "test", configs.Items[0].DataID)
	}
}

func TestListConfigInNs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testCsList))
	}))
	defer ts.Close()
	c := NewClient(ts.URL, "user", "password")
	configs, err := c.ListConfigInNs("test", "DEFAULT_GROUP")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(configs.Items))
		assert.Equal(t, "test", configs.Items[0].DataID)
	}
}

func TestListAllConfig(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.URL.Path == "/nacos/v1/console/namespaces" {
			w.Write([]byte(testNsList))
			return
		}
		w.Write([]byte(testCsList))
	}))
	defer ts.Close()
	c := NewClient(ts.URL, "user", "password")
	configs, err := c.ListAllConfig()
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(configs.Items))
		assert.Equal(t, "test", configs.Items[0].DataID)
	}
}

func TestCreateConfig(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	c.Token = &Token{AccessToken: "test-token"}
	err := c.CreateConfig(&CreateCSOpts{DataID: "test", Group: "DEFAULT_GROUP", Content: "test content", Tenant: "test-tenant", Type: "properties"})
	assert.NoError(t, err)
}

func TestDeleteConfig(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	c.Token = &Token{AccessToken: "test-token"}
	err := c.DeleteConfig(&DeleteCSOpts{DataID: "test", Group: "DEFAULT_GROUP", Tenant: "test-tenant"})
	assert.NoError(t, err)
}
