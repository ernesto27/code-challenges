package main

import (
	"reflect"
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
		{
			name:        "empty array",
			input:       `{"key": []}`,
			shouldError: false,
			expected:    map[string]interface{}{"key": []interface{}{}},
		},
		{
			name:        "array with numbers",
			input:       `{"numbers": [1, 2, 3]}`,
			shouldError: false,
			expected:    map[string]interface{}{"numbers": []interface{}{"1", "2", "3"}},
		},
		{
			name:        "array with mixed types",
			input:       `{"mixed": [1, "text", true, null]}`,
			shouldError: false,
			expected:    map[string]interface{}{"mixed": []interface{}{"1", "text", true, nil}},
		},
		{
			name:        "nested empty object",
			input:       `{"nested": {}}`,
			shouldError: false,
			expected:    map[string]interface{}{"nested": map[string]interface{}{}},
		},
		{
			name:        "nested object with values",
			input:       `{"outer": {"inner": "value"}}`,
			shouldError: false,
			expected:    map[string]interface{}{"outer": map[string]interface{}{"inner": "value"}},
		},
		{
			name:        "array with objects",
			input:       `{"objects": [{"id": 1}, {"id": 2}]}`,
			shouldError: false,
			expected:    map[string]interface{}{"objects": []interface{}{map[string]interface{}{"id": "1"}, map[string]interface{}{"id": "2"}}},
		},
		{
			name:        "complex nested structure",
			input:       `{"data": {"items": [{"name": "item1", "tags": ["tag1", "tag2"]}, {"name": "item2", "tags": []}]}}`,
			shouldError: false,
			expected:    map[string]interface{}{"data": map[string]interface{}{"items": []interface{}{map[string]interface{}{"name": "item1", "tags": []interface{}{"tag1", "tag2"}}, map[string]interface{}{"name": "item2", "tags": []interface{}{}}}}},
		},
		{
			name:        "unclosed array",
			input:       `{"key": [1, 2}`,
			shouldError: true,
			expected:    nil,
		},
		{
			name:        "missing comma in array",
			input:       `{"key": [1 2]}`,
			shouldError: true,
			expected:    nil,
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
				if !reflect.DeepEqual(result, tc.expected) {
					t.Errorf("Expected %+v, got %+v", tc.expected, result)
				}
			}
		})
	}
}
