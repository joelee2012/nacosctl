package nacos

import (
	"encoding/json"
	"io"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/jedib0t/go-pretty/table"
)

func toJson(v any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func toYaml(v any, w io.Writer) error {
	enc := yaml.NewEncoder(w)
	return enc.Encode(v)
}

func readYamlFile(v any, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := yaml.NewDecoder(f, yaml.DisallowUnknownField())
	return dec.Decode(v)
}

func writeYamlFile(v any, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return toYaml(v, f)
}

func toTable(w io.Writer, fn func(t table.Writer)) {
	tb := table.NewWriter()
	tb.SetOutputMirror(w)
	fn(tb)
	s := table.StyleLight
	s.Options = table.OptionsNoBordersAndSeparators
	tb.SetStyle(s)
	tb.Render()
}

func WriteAsFormat(format string, writable FormatWriter, w io.Writer) {
	switch format {
	case "json":
		writable.ToJson(w)
	case "yaml":
		writable.ToYaml(w)
	case "table":
		writable.ToTable(w)
	default:
		writable.ToTable(w)
	}
}

func WriteToDir(name string, writable DirWriter) error {
	return writable.WriteToDir(name)
}
func LoadFromYaml(name string, loader YamlFileLoader) error {
	return loader.FromYaml(name)
}
