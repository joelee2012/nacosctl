package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewNacos(t *testing.T) {
	n := NewNacos("http://localhost:8848", "user", "password")
	if n.URL != "http://localhost:8848" || n.User != "user" || n.Password != "password" {
		t.Errorf("NewNacos() failed, got: %v, want: %v", n, &Nacos{URL: "http://localhost:8848", User: "user", Password: "password"})
	}
}

func TestGetVersion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version": "1.0.0"}`))
	}))
	defer ts.Close()

	n := NewNacos(ts.URL, "user", "password")
	version, err := n.GetVersion()
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

	n := NewNacos(ts.URL, "user", "password")
	token, err := n.GetToken()
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

	n := NewNacos(ts.URL, "user", "password")
	namespaces, err := n.ListNamespace()
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

	n := NewNacos(ts.URL, "user", "password")
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

	n := NewNacos(ts.URL, "user", "password")
	err := n.DeleteNamespace("test-id")
	if err != nil {
		t.Errorf("DeleteNamespace() failed with error: %v", err)
	}
}

func TestUpdateNamespace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewNacos(ts.URL, "user", "password")
	err := n.UpdateNamespace(&CreateNSOpts{Name: "test", Desc: "Test namespace", ID: "test-id"})
	if err != nil {
		t.Errorf("UpdateNamespace() failed with error: %v", err)
	}
}

func TestListConfig(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"totalCount": 1, "pageNumber": 1, "pagesAvailable": 1, "pageItems": [{"id": "1", "dataId": "test", "group": "DEFAULT_GROUP", "content": "test content", "md5": "test-md5", "encryptedDataKey": "test-key", "tenant": "test-tenant", "appName": "test-app", "type": "properties"}]}`))
	}))
	defer ts.Close()

	n := NewNacos(ts.URL, "user", "password")
	configs, err := n.ListConfig(&ListCSOpts{DataId: "test", Group: "DEFAULT_GROUP", PageNumber: 1, PageSize: 10})
	if err != nil {
		t.Errorf("ListConfig() failed with error: %v", err)
	}
	if len(configs.PageItems) != 1 || configs.PageItems[0].DataID != "test" {
		t.Errorf("ListConfig() failed, got: %v, want: %v", configs, &ConfigList{PageItems: []*Config{{ID: "1", DataID: "test", Group: "DEFAULT_GROUP", Content: "test content", Md5: "test-md5", EncryptedDataKey: "test-key", Tenant: "test-tenant", AppName: "test-app", Type: "properties"}}})
	}
}

func TestCreateConfig(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewNacos(ts.URL, "user", "password")
	err := n.CreateConfig(&CreateCSOpts{DataID: "test", Group: "DEFAULT_GROUP", Content: "test content", Tenant: "test-tenant", Type: "properties"})
	if err != nil {
		t.Errorf("CreateConfig() failed with error: %v", err)
	}
}

func TestDeleteConfig(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewNacos(ts.URL, "user", "password")
	err := n.DeleteConfig(&CreateCSOpts{DataID: "test", Group: "DEFAULT_GROUP", Tenant: "test-tenant"})
	if err != nil {
		t.Errorf("DeleteConfig() failed with error: %v", err)
	}
}
