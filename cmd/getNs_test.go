/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import "testing"

func TestGetNamespace(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetNamespace(tt.args.args)
		})
	}
}
