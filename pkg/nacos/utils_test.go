package nacos

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/jedib0t/go-pretty/table"
	"github.com/stretchr/testify/assert"
)

func TestWriteJson(t *testing.T) {
	t.Run("successful write", func(t *testing.T) {
		var buf bytes.Buffer
		data := map[string]string{"key": "value"}
		err := writeJson(data, &buf)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), `"key": "value"`)
	})

	t.Run("write error", func(t *testing.T) {
		err := writeJson("data", errorWriter{})
		assert.Error(t, err)
	})
}

// TestWriteYaml tests the writeYaml function
func TestWriteYaml(t *testing.T) {
	t.Run("successful write", func(t *testing.T) {
		var buf bytes.Buffer
		data := map[string]string{"key": "value"}
		err := writeYaml(data, &buf)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "key: value")
	})

	// t.Run("write error", func(t *testing.T) {
	// 	data := map[string]string{"key": "value"}
	// 	err := writeYaml(data, errorWriter{})
	// 	assert.Error(t, err)
	// })
}

// TestWriteFile tests the writeFile function
func TestWriteFile(t *testing.T) {
	t.Run("successful write", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test.yaml")
		data := map[string]string{"key": "value"}
		err := writeYamlFile(data, tmpFile)
		assert.NoError(t, err)

		content, err := os.ReadFile(tmpFile)
		assert.NoError(t, err)
		assert.Contains(t, string(content), "key: value")
	})

	t.Run("file creation error", func(t *testing.T) {
		err := writeYamlFile("data", "/invalid/path/test.yaml")
		assert.Error(t, err)
	})
}

// TestWriteTable tests the writeTable function
func TestWriteTable(t *testing.T) {
	t.Run("table rendering", func(t *testing.T) {
		var buf bytes.Buffer
		writeTable(&buf, func(t table.Writer) {
			t.AppendHeader(table.Row{"Header"})
			t.AppendRow(table.Row{"Value"})
		})
		assert.Contains(t, buf.String(), "HEADER")
		assert.Contains(t, buf.String(), "Value")
	})
}
