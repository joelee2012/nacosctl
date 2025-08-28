package nacos

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		ns1File := filepath.Join(tmpDir, "ns1.yaml")
		assert.FileExists(t, ns1File)
	})

	t.Run("directory creation error", func(t *testing.T) {
		err := nl.WriteToDir("/invalid/path")
		assert.Error(t, err)
	})
}
