package main

import (
	"slices"
	"testing"
)

func TestSearchAll(t *testing.T) {

	inputs := []string{"echo", "exit", "type", "pwd"}
	trie := NewTrie()
	for _, input := range inputs {
		trie.Insert(input)
	}

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Basic Simple Test",
			input:    "e",
			expected: []string{"echo", "exit"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := trie.SearchAll(tt.input)

			if len(actual) != len(tt.expected) {
				t.Errorf("SearchAll(%s) = %v, expected: %v", tt.input, actual, tt.expected)
			}

			if !slices.Equal(actual, tt.expected) {
				t.Errorf("SearchAll(%s) = %v, expected: %v", tt.input, actual, tt.expected)
			}

		})
	}

}

func TestLcp(t *testing.T) {

	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "Basic Simple Test",
			input:    []string{"flower", "flow", "flight"},
			expected: "fl",
		},
		{
			name:     "Basic Simple Test",
			input:    []string{"flow", "flowhfdjasdhfa", "flowight"},
			expected: "flow",
		},
		{
			name:     "Basic Simple Test",
			input:    []string{"xyz_cow", "xyz_cow_pig", "xyz_cow_pig_bee"},
			expected: "xyz_cow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := NewTrie()
			trie.InsertAll(tt.input...)
			actual := trie.LongestCommonPrefix()

			if actual != tt.expected {
				t.Errorf("LongestCommonPrefix() = %v, expected: %v", actual, tt.expected)
			}

		})
	}

}
