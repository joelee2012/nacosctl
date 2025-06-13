package nacos

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("http://localhost:8848", "user", "password")
	if c.URL != "http://localhost:8848" || c.User != "user" || c.Password != "password" {
		t.Errorf("NewClient() failed, got: %v, want: %v", c, &Client{URL: "http://localhost:8848", User: "user", Password: "password"})
	}
}

func TestGetVersion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version": "1.0.0"}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	version, err := c.GetVersion()
	if err != nil {
		t.Errorf("GetVersion() failed with error: %v", err)
	}
	if version != "1.0.0" {
		t.Errorf("GetVersion() failed, got: %s, want: %s", version, "1.0.0")
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
	if err != nil {
		t.Errorf("GetToken() failed with error: %v", err)
	}
	if token != "test-token" {
		t.Errorf("GetToken() failed, got: %s, want: %s", token, "test-token")
	}
}

func TestListNamespace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 200, "message": "success", "data": [{"namespace": "test", "namespaceShowName": "Test", "namespaceDesc": "Test namespace", "quota": 100, "configCount": 10, "type": 0}]}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	namespaces, err := c.ListNamespace()
	if err != nil {
		t.Errorf("ListNamespace() failed with error: %v", err)
	}
	if len(namespaces.Items) != 1 || namespaces.Items[0].Name != "test" {
		t.Errorf("ListNamespace() failed, got: %v, want: %v", namespaces, &NsList{Items: []*Namespace{{Name: "test", ShowName: "Test", Desc: "Test namespace", Quota: 100, ConfigCount: 10, Type: 0}}})
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
	if err != nil {
		t.Errorf("CreateNamespace() failed with error: %v", err)
	}
}

func TestDeleteNamespace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	c.Token = &Token{AccessToken: "test-token"}
	err := c.DeleteNamespace("test-id")
	if err != nil {
		t.Errorf("DeleteNamespace() failed with error: %v", err)
	}
}

func TestUpdateNamespace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	c.Token = &Token{AccessToken: "test-token"}
	err := c.UpdateNamespace(&CreateNSOpts{Name: "test", Desc: "Test namespace", ID: "test-id"})
	if err != nil {
		t.Errorf("UpdateNamespace() failed with error: %v", err)
	}
}
func TestCreateOrUpdateNamespace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "success", "data": [{"namespace": "test", "namespaceShowName": "Test", "namespaceDesc": "Test namespace", "quota": 100, "configCount": 10, "type": 0}]}`))
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
		{name: "CreateNamespace", data: CreateNSOpts{Name: "test", Desc: "Test namespace", ID: "test"}},
		{name: "UpdateNamespace", data: CreateNSOpts{Name: "test-id", Desc: "Test namespace", ID: "test-id"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.CreateOrUpdateNamespace(&tt.data)
			if err != nil {
				t.Errorf("CreateOrUpdateNamespace() failed with error: %v", err)
			}
		})
	}
}

func TestListConfig(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"totalCount": 1, "pageNumber": 1, "pagesAvailable": 1, "pageItems": [{"id": "1", "dataId": "test", "group": "DEFAULT_GROUP", "content": "test content", "md5": "test-md5", "encryptedDataKey": "test-key", "tenant": "test-tenant", "appName": "test-app", "type": "properties"}]}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	configs, err := c.ListConfig(&ListCSOpts{DataId: "test", Group: "DEFAULT_GROUP", PageNumber: 1, PageSize: 10})
	if err != nil {
		t.Errorf("ListConfig() failed with error: %v", err)
	}
	if len(configs.Items) != 1 || configs.Items[0].DataID != "test" {
		t.Errorf("ListConfig() failed, got: %v, want: %v", configs, &ConfigList{Items: []*Config{{ID: "1", DataID: "test", Group: "DEFAULT_GROUP", Content: "test content", Md5: "test-md5", EncryptedDataKey: "test-key", Tenant: "test-tenant", AppName: "test-app", Type: "properties"}}})
	}
}

func TestListConfigInNs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"totalCount": 1, "pageNumber": 1, "pagesAvailable": 0, "pageItems": [{"id": "1", "dataId": "test", "group": "DEFAULT_GROUP", "content": "test content", "md5": "test-md5", "encryptedDataKey": "test-key", "tenant": "test-tenant", "appName": "test-app", "type": "properties"}]}`))
	}))
	defer ts.Close()
	c := NewClient(ts.URL, "user", "password")
	configs, err := c.ListConfigInNs("test", "DEFAULT_GROUP")
	if err != nil {
		t.Errorf("ListConfigInNs() failed with error: %v", err)
	}
	if len(configs.Items) != 1 || configs.Items[0].DataID != "test" {
		t.Errorf("ListConfigInNs() failed, got: %v, want: %v", configs, &ConfigList{Items: []*Config{{ID: "1", DataID: "test", Group: "DEFAULT_GROUP", Content: "test content", Md5: "test-md5", EncryptedDataKey: "test-key", Tenant: "test-tenant", AppName: "test-app", Type: "properties"}}})
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
	if err != nil {
		t.Errorf("CreateConfig() failed with error: %v", err)
	}
}

func TestDeleteConfig(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "password")
	c.Token = &Token{AccessToken: "test-token"}
	err := c.DeleteConfig(&DeleteCSOpts{DataID: "test", Group: "DEFAULT_GROUP", Tenant: "test-tenant"})
	if err != nil {
		t.Errorf("DeleteConfig() failed with error: %v", err)
	}
}
