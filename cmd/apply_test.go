package cmd

import (
	"testing"

	"github.com/joelee2012/nacosctl/pkg/nacos"
)

func TestCreateResourceFromFile(t *testing.T) {
	type args struct {
		naClient *nacos.Client
		name     string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateResourceFromFile(tt.args.naClient, tt.args.name)
		})
	}
}

func TestCreateResourceFromDir(t *testing.T) {
	type args struct {
		naClient *nacos.Client
		dir      string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateResourceFromDir(tt.args.naClient, tt.args.dir)
		})
	}
}
