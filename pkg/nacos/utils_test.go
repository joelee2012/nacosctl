package nacos

import (
	"bytes"
	"io"
	"path/filepath"
	"testing"

	"github.com/jedib0t/go-pretty/table"
	"github.com/stretchr/testify/assert"
)

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
			err := writeJson("data", tt.writer)
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
	err := writeYaml("data", &bytes.Buffer{})
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
	writeTable(&buf, func(t table.Writer) {
		t.AppendHeader(table.Row{"Header"})
		t.AppendRow(table.Row{"Value"})
	})
	assert.Contains(t, buf.String(), "HEADER")
	assert.Contains(t, buf.String(), "Value")
}

type mockFormatWriter map[string]bool

func (mw mockFormatWriter) WriteJson(w io.Writer) error {
	mw["json"] = true
	return nil
}
func (mw mockFormatWriter) WriteTable(w io.Writer) {
	mw["table"] = true
}

func (mw mockFormatWriter) WriteYaml(w io.Writer) error {
	mw["yaml"] = true
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
