package main

import (
	"testing"
)

func TestParseEmptyObject(t *testing.T) {
	input := "{}"

	lexer := NewLexer(input)
	parser := NewParser(lexer)

	result, err := parser.Parse()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty object, got: %v", result)
	}
}

func TestParseJSON(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		shouldError bool
		expected    map[string]interface{}
	}{
		{
			name:        "valid key-value",
			input:       `{"key": "value"}`,
			shouldError: false,
			expected:    map[string]interface{}{"key": "value"},
		},
		{
			name:        "missing colon",
			input:       `{"key" "value"}`,
			shouldError: true,
			expected:    nil,
		},
		{
			name:        "missing closing brace",
			input:       `{"key": "value"`,
			shouldError: true,
			expected:    nil,
		},
		{
			name:        "unquoted key",
			input:       `{key: "value"}`,
			shouldError: true,
			expected:    nil,
		},
		{
			name:        "empty key",
			input:       `{"": "value"}`,
			shouldError: false,
			expected:    map[string]interface{}{"": "value"},
		},
		{
			name:        "null value",
			input:       `{"key": null}`,
			shouldError: false,
			expected:    map[string]interface{}{"key": nil},
		},
		{
			name:        "true value",
			input:       `{"key": true}`,
			shouldError: false,
			expected:    map[string]interface{}{"key": true},
		},
		{
			name:        "number value",
			input:       `{"key": 111}`,
			shouldError: false,
			expected:    map[string]interface{}{"key": "111"},
		},
		{
			name:        "multi-key object",
			input:       `{"key1": true, "key2": false, "key3": null, "key4": "value", "key5": 101}`,
			shouldError: false,
			expected:    map[string]interface{}{"key1": true, "key2": false, "key3": nil, "key4": "value", "key5": "101"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lexer := NewLexer(tc.input)
			parser := NewParser(lexer)

			result, err := parser.Parse()

			if tc.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if result != nil {
					t.Errorf("Expected nil result but got: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if len(result) != len(tc.expected) {
					t.Errorf("Expected %d keys, got %d", len(tc.expected), len(result))
				}
				for key, expectedValue := range tc.expected {
					if actualValue, exists := result[key]; !exists {
						t.Errorf("Expected key '%s' to exist", key)
					} else if actualValue != expectedValue {
						t.Errorf("Expected value '%v', got '%v'", expectedValue, actualValue)
					}
				}
			}
		})
	}
}
