package nacos

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewNacos(t *testing.T) {
	n := NewClient("http://localhost:8848", "user", "password")
	if n.URL != "http://localhost:8848" || n.User != "user" || n.Password != "password" {
		t.Errorf("NewNacos() failed, got: %v, want: %v", n, &Client{URL: "http://localhost:8848", User: "user", Password: "password"})
	}
}

func TestGetVersion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version": "1.0.0"}`))
	}))
	defer ts.Close()

	n := NewClient(ts.URL, "user", "password")
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

	n := NewClient(ts.URL, "user", "password")
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

	n := NewClient(ts.URL, "user", "password")
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

	n := NewClient(ts.URL, "user", "password")
	n.Token = &Token{AccessToken: "test-token"}
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

	n := NewClient(ts.URL, "user", "password")
	n.Token = &Token{AccessToken: "test-token"}
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

	n := NewClient(ts.URL, "user", "password")
	configs, err := n.ListConfig(&ListCSOpts{DataId: "test", Group: "DEFAULT_GROUP", PageNumber: 1, PageSize: 10})
	if err != nil {
		t.Errorf("ListConfig() failed with error: %v", err)
	}
	if len(configs.Items) != 1 || configs.Items[0].DataID != "test" {
		t.Errorf("ListConfig() failed, got: %v, want: %v", configs, &ConfigList{Items: []*Config{{ID: "1", DataID: "test", Group: "DEFAULT_GROUP", Content: "test content", Md5: "test-md5", EncryptedDataKey: "test-key", Tenant: "test-tenant", AppName: "test-app", Type: "properties"}}})
	}
}

func TestCreateConfig(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewClient(ts.URL, "user", "password")
	n.Token = &Token{AccessToken: "test-token"}
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

	n := NewClient(ts.URL, "user", "password")
	n.Token = &Token{AccessToken: "test-token"}
	err := n.DeleteConfig(&DeleteCSOpts{DataID: "test", Group: "DEFAULT_GROUP", Tenant: "test-tenant"})
	if err != nil {
		t.Errorf("DeleteConfig() failed with error: %v", err)
	}
}

func TestNacos_GetVersion(t *testing.T) {
	tests := []struct {
		name    string
		n       *Client
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.GetVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("Nacos.GetVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Nacos.GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNacos_GetToken(t *testing.T) {
	tests := []struct {
		name    string
		n       *Client
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.GetToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("Nacos.GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Nacos.GetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNacos_ListNamespace(t *testing.T) {
	tests := []struct {
		name    string
		n       *Client
		want    *NsList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.ListNamespace()
			if (err != nil) != tt.wantErr {
				t.Errorf("Nacos.ListNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nacos.ListNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNacos_CreateNamespace(t *testing.T) {
	type args struct {
		opts *CreateNSOpts
	}
	tests := []struct {
		name    string
		n       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.CreateNamespace(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Nacos.CreateNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNacos_DeleteNamespace(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		n       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.DeleteNamespace(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Nacos.DeleteNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNacos_UpdateNamespace(t *testing.T) {
	type args struct {
		opts *CreateNSOpts
	}
	tests := []struct {
		name    string
		n       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.UpdateNamespace(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Nacos.UpdateNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNacos_CreateOrUpdateNamespace(t *testing.T) {
	type args struct {
		opts *CreateNSOpts
	}
	tests := []struct {
		name    string
		n       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.CreateOrUpdateNamespace(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Nacos.CreateOrUpdateNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNacos_ListConfig(t *testing.T) {
	type args struct {
		opts *ListCSOpts
	}
	tests := []struct {
		name    string
		n       *Client
		args    args
		want    *ConfigList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.ListConfig(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Nacos.ListConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nacos.ListConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNacos_ListConfigInNs(t *testing.T) {
	type args struct {
		namespace string
		group     string
	}
	tests := []struct {
		name    string
		n       *Client
		args    args
		want    *ConfigList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.ListConfigInNs(tt.args.namespace, tt.args.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("Nacos.ListConfigInNs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nacos.ListConfigInNs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNacos_ListAllConfig(t *testing.T) {
	tests := []struct {
		name    string
		n       *Client
		want    *ConfigList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.ListAllConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("Nacos.ListAllConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nacos.ListAllConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNacos_CreateConfig(t *testing.T) {
	type args struct {
		opts *CreateCSOpts
	}
	tests := []struct {
		name    string
		n       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.CreateConfig(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Nacos.CreateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNacos_DeleteConfig(t *testing.T) {
	type args struct {
		opts *DeleteCSOpts
	}
	tests := []struct {
		name    string
		n       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.DeleteConfig(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Nacos.DeleteConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
