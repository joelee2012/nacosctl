package cmd

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"testing"

	"github.com/jedib0t/go-pretty/table"
	"github.com/joelee2012/nacosctl/pkg/nacos"
	"github.com/stretchr/testify/assert"
)

// TestWriteJson tests the writeJson function

// TestConfigurationListWriteTable tests ConfigurationList.WriteTable method
func TestConfigurationListWriteTable(t *testing.T) {
	t.Run("with items", func(t *testing.T) {
		var buf bytes.Buffer
		cl := NewConfigurationList("v1", &nacos.ConfigList{
			Items: []*nacos.Config{
				{NamespaceID: "ns1", DataID: "data1", Group: "group1", Application: "app1", Type: "type1"},
				{NamespaceID: "ns2", DataID: "data2", Group: "group2", Application: "app2", Type: "type2"},
			},
		})
		cl.ToTable(&buf)
		output := buf.String()
		assert.Contains(t, output, "NAMESPACEID")
		assert.Contains(t, output, "ns1")
		assert.Contains(t, output, "data1")
	})

	t.Run("empty list", func(t *testing.T) {
		var buf bytes.Buffer
		cl := &ConfigurationList{Items: []*Configuration{}}
		cl.ToTable(&buf)
		assert.Equal(t, buf.String(), " NAMESPACEID  DATAID  GROUP  APPLICATION  TYPE \n")
	})
}

// TestConfigurationListWriteJson tests ConfigurationList.WriteJson method
func TestConfigurationListWriteJson(t *testing.T) {
	cl := NewConfigurationList("v1", &nacos.ConfigList{
		Items: []*nacos.Config{
			{NamespaceID: "ns1", DataID: "data1"},
		},
	})
	var buf bytes.Buffer
	err := cl.ToJson(&buf)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), `"namespace": "ns1"`)
}

// TestConfigurationListWriteToDir tests ConfigurationList.WriteToDir method
func TestConfigurationListWriteToDir(t *testing.T) {
	cl := NewConfigurationList("v1", &nacos.ConfigList{
		Items: []*nacos.Config{
			{NamespaceID: "ns1", DataID: "data1", Group: "group1"},
			{NamespaceID: "", DataID: "public_data", Group: "public_group"},
		},
	})
	t.Run("successful write", func(t *testing.T) {
		tmpDir := t.TempDir()
		err := cl.WriteToDir(tmpDir)
		assert.NoError(t, err)
		ns1File := filepath.Join(tmpDir, "ns1", "group1", "data1")
		assert.FileExists(t, ns1File)

		publicFile := filepath.Join(tmpDir, "public", "public_group", "public_data")
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
	tmpFile := filepath.Join(tmpDir, "config.yaml")
	config := NewConfiguration("v1", &nacos.Config{NamespaceID: "test"})
	err := config.ToFile(tmpFile)
	assert.NoError(t, err)
	assert.FileExists(t, tmpFile)
}

func TestNamespaceListWriteTable(t *testing.T) {
	t.Run("with items", func(t *testing.T) {
		var buf bytes.Buffer
		cl := NewNamespaceList("v1", &nacos.NsList{
			Items: []*nacos.Namespace{
				{ID: "ns1", Name: "data1", Description: "group1"},
			},
		})
		cl.ToTable(&buf)
		output := buf.String()
		assert.Contains(t, output, "NAME")
		assert.Contains(t, output, "ns1")
		assert.Contains(t, output, "data1")
	})

	t.Run("empty list", func(t *testing.T) {
		var buf bytes.Buffer
		cl := NewNamespaceList("v1", &nacos.NsList{Items: []*nacos.Namespace{}})
		cl.ToTable(&buf)
		assert.Equal(t, buf.String(), " NAME  ID  DESCRIPTION  COUNT \n")
	})
}

func TestNamespaceListWriteToDir(t *testing.T) {
	nl := NewNamespaceList("v1", &nacos.NsList{
		Items: []*nacos.Namespace{
			{ID: "ns1", Name: "showname", Description: "description"},
		},
	})
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

// TestWriteTable tests the writeTable function
func TestWriteTable(t *testing.T) {
	var buf bytes.Buffer
	toTable(&buf, func(t table.Writer) {
		t.AppendHeader(table.Row{"Header"})
		t.AppendRow(table.Row{"Value"})
	})
	assert.Contains(t, buf.String(), "HEADER")
	assert.Contains(t, buf.String(), "Value")
}

type mockFormatWriter map[string]bool

func (mw mockFormatWriter) ToJson(w io.Writer) error {
	mw["json"] = true
	return nil
}
func (mw mockFormatWriter) ToTable(w io.Writer) {
	mw["table"] = true
}

func (mw mockFormatWriter) ToYaml(w io.Writer) error {
	mw["yaml"] = true
	return nil
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
		{"json", true},
		{"yaml", true},
		{"table", true},
		{"xxx", false},
	}
	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			writer := mockFormatWriter{}
			WriteAsFormat(tt.format, writer, &bytes.Buffer{})
			assert.Equal(t, tt.called, writer[tt.format])
		})
	}
}
