package nacos

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// // errorWriter is a mock io.Writer that always returns an error
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
