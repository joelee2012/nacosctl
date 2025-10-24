package cmd

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"testing"

	"github.com/joelee2012/nacosctl/pkg/nacos"
	"github.com/stretchr/testify/assert"
)

var apiVersion = "v1"
var cs = []nacos.Configuration{
	{NamespaceID: "ns1", DataID: "data1", Group: "group1", Application: "app1", Type: "type1"},
	{NamespaceID: "ns2", DataID: "data2", Group: "group2", Application: "app2", Type: "type2"},
}
var ns = []nacos.Namespace{
	{ID: "ns1", Name: "data1", Description: "group1"},
}

var us = []nacos.User{
	{Name: "user1", Password: "password1"},
}

// TestConfigurationListWriteTable tests ConfigurationList.WriteTable method
func TestConfigurationListWriteTable(t *testing.T) {
	t.Run("with items", func(t *testing.T) {
		var buf bytes.Buffer
		cl := NewList(apiVersion, cs, NewConfiguration)
		cl.ToTable(&buf)
		output := buf.String()
		assert.Contains(t, output, "NAMESPACEID")
		assert.Contains(t, output, "ns1")
		assert.Contains(t, output, "data1")
	})

	t.Run("empty list", func(t *testing.T) {
		var buf bytes.Buffer
		cl := List[Configuration]{}
		cl.ToTable(&buf)
		assert.Equal(t, "No resources found", buf.String())
	})
}

// TestConfigurationListWriteToDir tests ConfigurationList.WriteToDir method
func TestConfigurationListWriteToDir(t *testing.T) {
	cl := NewList(apiVersion, cs, NewConfiguration)
	t.Run("successful write", func(t *testing.T) {
		tmpDir := t.TempDir()
		err := cl.WriteToDir(tmpDir)
		assert.NoError(t, err)
		ns1File := filepath.Join(tmpDir, "ns1", "group1", "data1")
		assert.FileExists(t, ns1File)

		publicFile := filepath.Join(tmpDir, "ns2", "group2", "data2")
		assert.FileExists(t, publicFile)
	})

	t.Run("directory creation error", func(t *testing.T) {
		err := cl.WriteToDir("/invalid/path")
		assert.Error(t, err)
	})
}

// TestConfigWriteFile tests Config.WriteFile method
func TestConfigWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	c := NewConfiguration(apiVersion, cs[0])
	tmpFile := filepath.Join(tmpDir, c.Metadata.Namespace, c.Metadata.Group, c.Metadata.DataID)
	err := c.WriteToDir(tmpDir)
	assert.NoError(t, err)
	assert.FileExists(t, tmpFile)
}

func TestNamespaceListWriteTable(t *testing.T) {
	t.Run("with items", func(t *testing.T) {
		var buf bytes.Buffer
		cl := NewList(apiVersion, ns, NewNamespace)
		cl.ToTable(&buf)
		output := buf.String()
		assert.Contains(t, output, "NAME")
		assert.Contains(t, output, "ns1")
		assert.Contains(t, output, "data1")
	})

	t.Run("empty list", func(t *testing.T) {
		var buf bytes.Buffer
		cl := List[Namespace]{}
		cl.ToTable(&buf)
		assert.Equal(t, "No resources found", buf.String())
	})
}

func TestNamespaceListWriteToDir(t *testing.T) {
	nl := NewList(apiVersion, ns, NewNamespace)
	t.Run("successful write", func(t *testing.T) {
		tmpDir := t.TempDir()
		err := nl.WriteToDir(tmpDir)
		assert.NoError(t, err)
		ns1File := filepath.Join(tmpDir, "ns1.yaml")
		assert.FileExists(t, ns1File)
	})

	t.Run("directory creation error", func(t *testing.T) {
		err := nl.WriteToDir("/invalid/path")
		assert.Error(t, err)
	})
}

type errorWriter struct{}

func (ew *errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("mock write error")
}

func TestWriteJson(t *testing.T) {
	tests := []struct {
		name    string
		writer  io.Writer
		wantErr bool
	}{
		{"OK", &bytes.Buffer{}, false},
		{"NOK", &errorWriter{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := toJson("data", tt.writer)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestWriteYaml tests the writeYaml function
func TestWriteYaml(t *testing.T) {
	err := toYaml("data", &bytes.Buffer{})
	assert.NoError(t, err)
}

// TestWriteFile tests the writeFile function
func TestWriteYamlFile(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{"OK", filepath.Join(t.TempDir(), "test.yaml"), false},
		{"NOK", "/invalied/test.yaml", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := writeYamlFile("a", tt.file)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReadYamlFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.yaml")
	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{"OK", tmpFile, false},
		{"NOK", "/invalied/test.yaml", true},
	}
	err := writeYamlFile(map[string]int{"a": 1}, tmpFile)
	assert.NoError(t, err)
	var v map[string]int
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := readYamlFile(&v, tt.file)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, v["a"], 1)
			}
		})
	}
}

type mockFormatWriter map[string]bool

func (mw mockFormatWriter) ToTable(w io.Writer) {
	mw["table"] = true
}

func (mw mockFormatWriter) WriteToDir(w string) error {
	mw["dir"] = true
	return nil
}
func TestWriteAsFormat(t *testing.T) {
	tests := []struct {
		format string
		called bool
	}{
		{"table", true},
		{"xxx", false},
	}
	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			writer := mockFormatWriter{}
			WriteFormat(writer, tt.format, &bytes.Buffer{})
			assert.Equal(t, tt.called, writer[tt.format])
		})
	}
}

func TestUserListWriteTable(t *testing.T) {
	t.Run("with items", func(t *testing.T) {
		var buf bytes.Buffer
		ul := NewList(apiVersion, us, NewUser)
		ul.ToTable(&buf)
		output := buf.String()
		assert.Contains(t, output, "NAME")
		assert.Contains(t, output, "user1")
		assert.Contains(t, output, "password1")
	})

	t.Run("empty list", func(t *testing.T) {
		var buf bytes.Buffer
		cl := List[User]{}
		cl.ToTable(&buf)
		assert.Equal(t, "No resources found", buf.String())
	})
}
