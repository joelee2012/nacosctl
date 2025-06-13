package nacos

import (
	"encoding/json"
	"io"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/jedib0t/go-pretty/table"
)

func writeJson(v any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func writeYaml(v any, w io.Writer) error {
	enc := yaml.NewEncoder(w)
	return enc.Encode(v)
}

func readYamlFile(v any, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(v)
}

func writeYamlFile(v any, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return writeYaml(v, f)
}

func writeTable(w io.Writer, fn func(t table.Writer)) {
	tb := table.NewWriter()
	tb.SetOutputMirror(w)
	fn(tb)
	s := table.StyleLight
	s.Options = table.OptionsNoBordersAndSeparators
	tb.SetStyle(s)
	tb.Render()
}
func WriteAsFormat(format string, writable FormatWriter) {
	switch format {
	case "json":
		writable.WriteJson(os.Stdout)
	case "yaml":
		writable.WriteYaml(os.Stdout)
	case "table":
		writable.WriteTable(os.Stdout)
	default:
		writable.WriteTable(os.Stdout)
	}
}

func WriteToDir(name string, writable DirWriter) error {
	return writable.WriteToDir(name)
}
func LoadFromYaml(name string, loader YamlFileLoader) error {
	return loader.LoadFromYaml(name)
}
