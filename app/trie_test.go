package main

import (
	"slices"
	"testing"
)

func TestSearcAll(t *testing.T) {

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
