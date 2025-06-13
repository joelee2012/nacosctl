/*
Copyright © 2025 Joe Lee <lj_2005@163.com>
*/
package cmd

import (
	"reflect"
	"testing"

	"github.com/joelee2012/nacosctl/pkg/nacos"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Execute()
		})
	}
}

func Test_initConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initConfig()
		})
	}
}

func TestNewNacosClient(t *testing.T) {
	tests := []struct {
		name string
		want *nacos.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNacosClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNacosClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
