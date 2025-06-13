package nacos

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/jedib0t/go-pretty/table"
	"github.com/stretchr/testify/assert"
)

// TestWriteJson tests the writeJson function
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

// mockFormatWriter is a mock implementation of FormatWriter for testing
type mockFormatWriter struct {
	tableCalled bool
	jsonCalled  bool
	yamlCalled  bool
}

func (m *mockFormatWriter) WriteTable(w io.Writer) {
	m.tableCalled = true
}

func (m *mockFormatWriter) WriteJson(w io.Writer) error {
	m.jsonCalled = true
	return nil
}

func (m *mockFormatWriter) WriteYaml(w io.Writer) error {
	m.yamlCalled = true
	return nil
}

func Test_writeJson(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := writeJson(tt.args.v, w); (err != nil) != tt.wantErr {
				t.Errorf("writeJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeJson() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeYaml(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := writeYaml(tt.args.v, w); (err != nil) != tt.wantErr {
				t.Errorf("writeYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeYaml() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_writeFile(t *testing.T) {
	type args struct {
		v    any
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeYamlFile(tt.args.v, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("writeFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_writeTable(t *testing.T) {
	type args struct {
		fn func(t table.Writer)
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			writeTable(w, tt.args.fn)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeTable() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestConfigList_WriteTable(t *testing.T) {
	tests := []struct {
		name  string
		c     *ConfigList
		wantW string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.WriteTable(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ConfigList.WriteTable() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestConfigList_WriteJson(t *testing.T) {
	tests := []struct {
		name    string
		c       *ConfigList
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.c.WriteJson(w); (err != nil) != tt.wantErr {
				t.Errorf("ConfigList.WriteJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ConfigList.WriteJson() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestConfigList_WriteYaml(t *testing.T) {
	tests := []struct {
		name    string
		c       *ConfigList
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.c.WriteYaml(w); (err != nil) != tt.wantErr {
				t.Errorf("ConfigList.WriteYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ConfigList.WriteYaml() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestConfigList_WriteToDir(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		cs      *ConfigList
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cs.WriteToDir(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ConfigList.WriteToDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_WriteJson(t *testing.T) {
	tests := []struct {
		name    string
		c       *Config
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.c.WriteJson(w); (err != nil) != tt.wantErr {
				t.Errorf("Config.WriteJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Config.WriteJson() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestConfig_WriteYaml(t *testing.T) {
	tests := []struct {
		name    string
		c       *Config
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.c.WriteYaml(w); (err != nil) != tt.wantErr {
				t.Errorf("Config.WriteYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Config.WriteYaml() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestConfig_WriteFile(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		c       *Config
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.WriteFile(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Config.WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_readYaml(t *testing.T) {
	type args struct {
		v    any
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := readYamlFile(tt.args.v, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("readYaml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_LoadFromYaml(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		c       *Config
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.LoadFromYaml(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Config.LoadFromYaml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNsList_WriteTable(t *testing.T) {
	tests := []struct {
		name  string
		n     *NsList
		wantW string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.n.WriteTable(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("NsList.WriteTable() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestNsList_WriteJson(t *testing.T) {
	tests := []struct {
		name    string
		n       *NsList
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.n.WriteJson(w); (err != nil) != tt.wantErr {
				t.Errorf("NsList.WriteJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("NsList.WriteJson() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestNsList_WriteYaml(t *testing.T) {
	tests := []struct {
		name    string
		n       *NsList
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.n.WriteYaml(w); (err != nil) != tt.wantErr {
				t.Errorf("NsList.WriteYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("NsList.WriteYaml() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestNsList_WriteToDir(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		n       *NsList
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.WriteToDir(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("NsList.WriteToDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNamespace_WriteJson(t *testing.T) {
	tests := []struct {
		name    string
		n       *Namespace
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.n.WriteJson(w); (err != nil) != tt.wantErr {
				t.Errorf("Namespace.WriteJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Namespace.WriteJson() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestNamespace_WriteYaml(t *testing.T) {
	tests := []struct {
		name    string
		n       *Namespace
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.n.WriteYaml(w); (err != nil) != tt.wantErr {
				t.Errorf("Namespace.WriteYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Namespace.WriteYaml() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestNamespace_WriteFile(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		n       *Namespace
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.WriteFile(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Namespace.WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNamespace_LoadFromYaml(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		n       *Namespace
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.LoadFromYaml(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Namespace.LoadFromYaml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
