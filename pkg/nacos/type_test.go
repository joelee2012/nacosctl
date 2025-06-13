package nacos

import (
	"bytes"
	"errors"
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
		assert.Contains(t, buf.String(), "NAMESPACE")
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
	t.Run("successful write", func(t *testing.T) {
		tmpDir := t.TempDir()
		cl := &ConfigList{
			Items: []*Config{
				{Tenant: "ns1", DataID: "data1", Group: "group1"},
				{Tenant: "", DataID: "public_data", Group: "public_group"},
			},
		}
		err := cl.WriteToDir(tmpDir)
		assert.NoError(t, err)

		// Check files were created
		ns1File := filepath.Join(tmpDir, "ns1", "group1", "data1")
		assert.FileExists(t, ns1File)

		publicFile := filepath.Join(tmpDir, "public", "public_group", "public_data")
		assert.FileExists(t, publicFile)
	})

	t.Run("directory creation error", func(t *testing.T) {
		cl := &ConfigList{
			Items: []*Config{
				{Tenant: "ns1", DataID: "data1", Group: "group1"},
			},
		}
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

// // TestLoadFromYaml tests LoadFromYaml function
// func TestLoadFromYaml(t *testing.T) {
// 	t.Run("successful load", func(t *testing.T) {
// 		tmpDir := t.TempDir()
// 		tmpFile := filepath.Join(tmpDir, "test.yaml")
// 		err := os.WriteFile(tmpFile, []byte("Tenant: test"), 0644)
// 		require.NoError(t, err)

// 		var ns Namespace
// 		err = LoadFromYaml(tmpFile, &ns)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "test", ns.Tenant)
// 	})

// 	t.Run("file not found", func(t *testing.T) {
// 		var ns Namespace
// 		err := LoadFromYaml("nonexistent.yaml", &ns)
// 		assert.Error(t, err)
// 	})
// }

// // TestWriteAsFormat tests WriteAsFormat function
// func TestWriteAsFormat(t *testing.T) {
// 	t.Run("json format", func(t *testing.T) {
// 		var buf bytes.Buffer
// 		oldStdout := os.Stdout
// 		defer func() { os.Stdout = oldStdout }()
// 		os.Stdout = &buf

// 		mockWritable := &mockFormatWriter{}
// 		WriteAsFormat("json", mockWritable)
// 		assert.True(t, mockWritable.jsonCalled)
// 	})

// 	t.Run("default to table", func(t *testing.T) {
// 		var buf bytes.Buffer
// 		oldStdout := os.Stdout
// 		defer func() { os.Stdout = oldStdout }()
// 		os.Stdout = &buf

// 		mockWritable := &mockFormatWriter{}
// 		WriteAsFormat("invalid", mockWritable)
// 		assert.True(t, mockWritable.tableCalled)
// 	})
// }

// errorWriter is a mock io.Writer that always returns an error
type errorWriter struct{}

func (ew errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("mock write error")
}
