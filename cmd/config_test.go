package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestReadFile tests the ReadFile method of CLIConfig.
func TestReadFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		wantErr     bool
	}{
		{"ValidYAML", "context: test\nservers:\n  server1:\n    password: pass1\n    url: url1\n    user: user1", false},
		{"InvalidYAML", "invalid yaml content", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "nacos.yaml")
			if err := os.WriteFile(tmpFile, []byte(tt.fileContent), 0666); err != nil {
				t.Error(err)
			}
			c := &CLIConfig{}
			err := c.ReadFile(tmpFile)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

var config = &CLIConfig{
	Context: "server1",
	Servers: map[string]*Server{
		"server1": {Password: "pass1", URL: "url1", User: "user1"},
		"server3": {Password: "pass3", URL: "url3", User: "user3"},
	},
}

// TestGetServer tests the GetServer method of CLIConfig.
func TestGetServer(t *testing.T) {
	tests := []struct {
		name     string
		server   string
		expected *Server
	}{
		{"ServerExists", "server1", &Server{Password: "pass1", URL: "url1", User: "user1"}},
		{"ServerNotExists", "notexist", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.GetServer(tt.server)
			assert.Equal(t, tt.expected, result)
			if result != nil {
				assert.Equal(t, tt.expected.URL, result.URL)
				assert.Equal(t, tt.expected.Password, result.Password)
				assert.Equal(t, tt.expected.User, result.User)
			}
		})
	}
}

// TestAddServer tests the AddServer method of CLIConfig.
func TestAddServer(t *testing.T) {
	config := &CLIConfig{
		Servers: make(map[string]*Server),
	}

	server := &Server{Password: "pass2", URL: "url2", User: "user2"}
	config.AddServer("server2", server)
	assert.Equal(t, server, config.Servers["server2"])
}

// TestDeleteServer tests the DeleteServer method of CLIConfig.
func TestDeleteServer(t *testing.T) {
	c := &CLIConfig{
		Context: "server3",
		Servers: map[string]*Server{
			"server3": {Password: "pass3", URL: "url3", User: "user3"},
		},
	}
	c.DeleteServer("server3")
	exist := c.Servers["server3"]
	assert.Nil(t, exist)
	assert.Empty(t, c.Context)
}

// TestSetContext tests the SetContext method of CLIConfig.
func TestSetContext(t *testing.T) {
	tests := []struct {
		name    string
		context string
		wantErr bool
	}{
		{"ValidContext", "server1", false},
		{"InvalidContext", "notexist", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.SetContext(tt.context)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.context, config.Context)
			}
		})
	}
}

// TestGetContext tests the GetContext method of CLIConfig.
func TestGetContext(t *testing.T) {
	assert.Equal(t, "server1", config.GetContext())
}

// TestGetCurrentServer tests the GetCurrentServer method of CLIConfig.
func TestGetCurrentServer(t *testing.T) {
	server := config.GetCurrentServer()
	assert.Equal(t, "url1", server.URL)
}

// TestToYaml tests the ToYaml method of CLIConfig.
func TestToYaml(t *testing.T) {
	_, err := config.ToYaml()
	assert.NoError(t, err)
}
