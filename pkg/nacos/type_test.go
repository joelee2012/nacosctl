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
				{NamespaceID: "ns1", DataID: "data1", Group: "group1", Application: "app1", Type: "type1"},
				{NamespaceID: "ns2", DataID: "data2", Group: "group2", Application: "app2", Type: "type2"},
			},
		}
		cl.ToTable(&buf)
		output := buf.String()
		assert.Contains(t, output, "NAMESPACEID")
		assert.Contains(t, output, "ns1")
		assert.Contains(t, output, "data1")
	})

	t.Run("empty list", func(t *testing.T) {
		var buf bytes.Buffer
		cl := &ConfigList{Items: []*Config{}}
		cl.ToTable(&buf)
		assert.Equal(t, buf.String(), " NAMESPACEID  DATAID  GROUP  APPLICATION  TYPE \n")
	})
}

// TestConfigListWriteJson tests ConfigList.WriteJson method
func TestConfigListWriteJson(t *testing.T) {
	cl := &ConfigList{
		Items: []*Config{
			{NamespaceID: "ns1", DataID: "data1"},
		},
	}
	var buf bytes.Buffer
	err := cl.ToJson(&buf)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), `"tenant": "ns1"`)
}

// TestConfigListWriteToDir tests ConfigList.WriteToDir method
func TestConfigListWriteToDir(t *testing.T) {
	cl := &ConfigList{
		Items: []*Config{
			{NamespaceID: "ns1", DataID: "data1", Group: "group1"},
			{NamespaceID: "", DataID: "public_data", Group: "public_group"},
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
	config := &Config{NamespaceID: "test"}
	err := config.ToFile(tmpFile)
	assert.NoError(t, err)
	assert.FileExists(t, tmpFile)
}

func TestNamespaceListWriteTable(t *testing.T) {
	t.Run("with items", func(t *testing.T) {
		var buf bytes.Buffer
		cl := &NsList{
			Items: []*Namespace{
				{ID: "ns1", Name: "data1", Description: "group1"},
			},
		}
		cl.ToTable(&buf)
		output := buf.String()
		assert.Contains(t, output, "NAME")
		assert.Contains(t, output, "ns1")
		assert.Contains(t, output, "data1")
	})

	t.Run("empty list", func(t *testing.T) {
		var buf bytes.Buffer
		cl := &NsList{Items: []*Namespace{}}
		cl.ToTable(&buf)
		assert.Equal(t, buf.String(), " NAME  ID  DESCRIPTION  COUNT \n")
	})
}

func TestNamespaceListWriteToDir(t *testing.T) {
	nl := &NsList{
		Items: []*Namespace{
			{ID: "ns1", Name: "showname", Description: "description"},
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
