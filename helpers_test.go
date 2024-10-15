package main

import (
	"testing"
)

func TestReplaceBadWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"ReplacesKnownBadWords", "This is a kerfuffle and a sharbert", "This is a **** and a ****"},
		{"IgnoresUnknownWords", "This is a test", "This is a test"},
		{"HandlesMixedCase", "This is a Kerfuffle and a SHARBERT", "This is a **** and a ****"},
		{"HandlesEmptyString", "", ""},
		{"HandlesMultipleSpaces", "This  is   a    kerfuffle", "This  is   a    ****"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceBadWords(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s but got %s", tt.expected, result)
			}
		})
	}
}
