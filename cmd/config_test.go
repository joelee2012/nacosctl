package cmd

import (
	"os"
	"testing"
)

// MockFileSystem is a mock implementation of the file system operations.
type MockFileSystem struct {
	OpenFunc  func(name string) (*os.File, error)
	WriteFunc func(data []byte, perm os.FileMode) error
}

func (m *MockFileSystem) Open(name string) (*os.File, error) {
	return m.OpenFunc(name)
}

func (m *MockFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	return m.WriteFunc(data, perm)
}

// TestReadFile tests the ReadFile method of CLIConfig.
func TestReadFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		expectError bool
	}{
		{"ValidYAML", "context: test\nservers:\n  server1:\n    password: pass1\n    url: url1\n    user: user1", false},
		{"InvalidYAML", "invalid yaml content", true},
	}

	f, err := os.CreateTemp("", "nacosctl")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name()) // clean up
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.WriteFile(f.Name(), []byte(tt.fileContent), 0666); err != nil {
				t.Error(err)
			}
			config := &CLIConfig{}
			err := config.ReadFile(f.Name())
			if (err != nil) != tt.expectError {
				t.Errorf("ReadFile() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

// TestWriteFile tests the WriteFile method of CLIConfig.
// func TestWriteFile(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		expectError bool
// 	}{
// 		{"WriteSuccess", false},
// 		{"WriteFailure", true},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockFS := &MockFileSystem{
// 				WriteFunc: func(data []byte, perm os.FileMode) error {
// 					if tt.expectError {
// 						return errors.New("write error")
// 					}
// 					return nil
// 				},
// 			}

// 			config := &CLIConfig{}
// 			err := config.WriteFile("testfile")
// 			if (err != nil) != tt.expectError {
// 				t.Errorf("WriteFile() error = %v, expectError %v", err, tt.expectError)
// 			}
// 		})
// 	}
// }

// TestGetServer tests the GetServer method of CLIConfig.
func TestGetServer(t *testing.T) {
	config := &CLIConfig{
		Servers: map[string]*Server{
			"server1": {Password: "pass1", URL: "url1", User: "user1"},
		},
	}

	tests := []struct {
		name     string
		server   string
		expected *Server
	}{
		{"ServerExists", "server1", &Server{Password: "pass1", URL: "url1", User: "user1"}},
		{"ServerNotExists", "server2", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.GetServer(tt.server)
			if result != tt.expected && (result.URL != tt.expected.URL || result.Password != tt.expected.Password || result.User != tt.expected.User) {
				t.Errorf("GetServer() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestAddServer tests the AddServer method of CLIConfig.
func TestAddServer(t *testing.T) {
	config := &CLIConfig{
		Servers: make(map[string]*Server),
	}

	server := &Server{Password: "pass1", URL: "url1", User: "user1"}
	config.AddServer("server1", server)

	if config.Servers["server1"] != server {
		t.Errorf("AddServer() failed to add server")
	}
}

// TestDeleteServer tests the DeleteServer method of CLIConfig.
func TestDeleteServer(t *testing.T) {
	config := &CLIConfig{
		Context: "server1",
		Servers: map[string]*Server{
			"server1": {Password: "pass1", URL: "url1", User: "user1"},
		},
	}

	config.DeleteServer("server1")

	if _, exists := config.Servers["server1"]; exists {
		t.Errorf("DeleteServer() failed to delete server")
	}

	if config.Context != "" {
		t.Errorf("DeleteServer() failed to clear context")
	}
}

// TestSetContext tests the SetContext method of CLIConfig.
func TestSetContext(t *testing.T) {
	config := &CLIConfig{
		Servers: map[string]*Server{
			"server1": {Password: "pass1", URL: "url1", User: "user1"},
		},
	}

	tests := []struct {
		name        string
		context     string
		expectError bool
	}{
		{"ValidContext", "server1", false},
		{"InvalidContext", "server2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.SetContext(tt.context)
			if (err != nil) != tt.expectError {
				t.Errorf("SetContext() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

// TestGetContext tests the GetContext method of CLIConfig.
func TestGetContext(t *testing.T) {
	config := &CLIConfig{
		Context: "server1",
	}

	if config.GetContext() != "server1" {
		t.Errorf("GetContext() = %v, expected server1", config.GetContext())
	}
}

// TestGetCurrentServer tests the GetCurrentServer method of CLIConfig.
func TestGetCurrentServer(t *testing.T) {
	config := &CLIConfig{
		Context: "server1",
		Servers: map[string]*Server{
			"server1": {Password: "pass1", URL: "url1", User: "user1"},
		},
	}

	server := config.GetCurrentServer()
	if server != config.Servers["server1"] {
		t.Errorf("GetCurrentServer() = %v, expected %v", server, config.Servers["server1"])
	}
}

// TestToYaml tests the ToYaml method of CLIConfig.
func TestToYaml(t *testing.T) {
	config := &CLIConfig{
		Context: "server1",
		Servers: map[string]*Server{
			"server1": {Password: "pass1", URL: "url1", User: "user1"},
		},
	}

	_, err := config.ToYaml()
	if err != nil {
		t.Errorf("ToYaml() error = %v", err)
	}
}
