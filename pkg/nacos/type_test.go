package nacos

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWriteJson tests the writeJson function

// TestConfigListWriteTable tests ConfigList.WriteTable method
func TestConfigListWriteTable(t *testing.T) {
	t.Run("with items", func(t *testing.T) {
		var buf bytes.Buffer
		cl := &ConfigList{
			Items: []*Config{
				{Tenant: "ns1", DataID: "data1", Group: "group1", AppName: "app1", Type: "type1"},
				{Tenant: "ns2", DataID: "data2", Group: "group2", AppName: "app2", Type: "type2"},
			},
		}
		cl.WriteTable(&buf)
		output := buf.String()
		assert.Contains(t, output, "NAMESPACE")
		assert.Contains(t, output, "ns1")
		assert.Contains(t, output, "data1")
	})

	t.Run("empty list", func(t *testing.T) {
		var buf bytes.Buffer
		cl := &ConfigList{Items: []*Config{}}
		cl.WriteTable(&buf)
		assert.Equal(t, buf.String(), " NAMESPACE  DATAID  GROUP  APPLICATION  TYPE \n")
	})
}

// TestConfigListWriteJson tests ConfigList.WriteJson method
func TestConfigListWriteJson(t *testing.T) {
	cl := &ConfigList{
		Items: []*Config{
			{Tenant: "ns1", DataID: "data1"},
		},
	}
	var buf bytes.Buffer
	err := cl.WriteJson(&buf)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), `"tenant": "ns1"`)
}

// TestConfigListWriteToDir tests ConfigList.WriteToDir method
func TestConfigListWriteToDir(t *testing.T) {
	cl := &ConfigList{
		Items: []*Config{
			{Tenant: "ns1", DataID: "data1", Group: "group1"},
			{Tenant: "", DataID: "public_data", Group: "public_group"},
		},
	}
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
	config := &Config{Tenant: "test"}
	err := config.WriteFile(tmpFile)
	assert.NoError(t, err)
	assert.FileExists(t, tmpFile)
}

func TestNamespaceListWriteTable(t *testing.T) {
	t.Run("with items", func(t *testing.T) {
		var buf bytes.Buffer
		cl := &NsList{
			Items: []*Namespace{
				{Name: "ns1", ShowName: "data1", Desc: "group1"},
			},
		}
		cl.WriteTable(&buf)
		output := buf.String()
		assert.Contains(t, output, "NAME")
		assert.Contains(t, output, "ns1")
		assert.Contains(t, output, "data1")
	})

	t.Run("empty list", func(t *testing.T) {
		var buf bytes.Buffer
		cl := &NsList{Items: []*Namespace{}}
		cl.WriteTable(&buf)
		assert.Equal(t, buf.String(), " NAME  ID  DESCRIPTION  COUNT \n")
	})
}

func TestNamespaceListWriteToDir(t *testing.T) {
	nl := &NsList{
		Items: []*Namespace{
			{Name: "ns1", ShowName: "showname", Desc: "description"},
		},
	}
	t.Run("successful write", func(t *testing.T) {
		tmpDir := t.TempDir()
		err := nl.WriteToDir(tmpDir)
		assert.NoError(t, err)
		ns1File := filepath.Join(tmpDir, "showname.yaml")
		assert.FileExists(t, ns1File)
	})

	t.Run("directory creation error", func(t *testing.T) {
		err := nl.WriteToDir("/invalid/path")
		assert.Error(t, err)
	})
}
