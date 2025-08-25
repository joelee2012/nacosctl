package nacos

import (
	"encoding/json"
	"io"
	"os"

	"github.com/goccy/go-yaml"
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
