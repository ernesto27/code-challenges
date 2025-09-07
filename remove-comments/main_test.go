package main

import (
	"testing"
)

func TestIsInsideString(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		pos      int
		expected bool
	}{
		{
			name:     "Hash outside string",
			line:     `print("Hello")  # comment`,
			pos:      16,
			expected: false,
		},
		{
			name:     "Hash inside double quotes",
			line:     `print("Hello # world")`,
			pos:      13,
			expected: true,
		},
		{
			name:     "Hash inside single quotes",
			line:     `print('Hello # world')`,
			pos:      13,
			expected: true,
		},
		{
			name:     "Hash after string ends",
			line:     `x = "test" # comment`,
			pos:      11,
			expected: false,
		},
		{
			name:     "Hash inside triple double quotes",
			line:     `"""docstring # here"""`,
			pos:      13,
			expected: true,
		},
		{
			name:     "Hash inside triple single quotes",
			line:     `'''docstring # here'''`,
			pos:      13,
			expected: true,
		},
		{
			name:     "Hash inside string with escaped quotes",
			line:     `print("He said \"Hi\" # test")`,
			pos:      20,
			expected: true,
		},
		{
			name:     "Hash inside string with escaped quote",
			line:     `print('Don\'t # remove')`,
			pos:      11,
			expected: true,
		},
		{
			name:     "Hash at start of line",
			line:     `# start of line`,
			pos:      0,
			expected: false,
		},
		{
			name:     "Simple inline comment",
			line:     `x = 5  # inline comment`,
			pos:      7,
			expected: false,
		},
		{
			name:     "Hash inside double quotes with single quotes",
			line:     `print("Mix 'quotes' # here")`,
			pos:      20,
			expected: true,
		},
		{
			name:     "Hash inside single quotes with double quotes",
			line:     `print('Mix "quotes" # here')`,
			pos:      20,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rc := &RemoveComments{}
			result := rc.isInsideString(tc.line, tc.pos)
			if result != tc.expected {
				t.Errorf("isInsideString(%q, %d) = %v; expected %v",
					tc.line, tc.pos, result, tc.expected)
			}
		})
	}
}

func TestRemoveCommentsFromLine(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple comment removal",
			input:    `x = 5  # this is a comment`,
			expected: `x = 5`,
		},
		{
			name:     "No comment in line",
			input:    `x = 5`,
			expected: `x = 5`,
		},
		{
			name:     "Hash inside string should be preserved",
			input:    `print("Hello # World")`,
			expected: `print("Hello # World")`,
		},
		{
			name:     "Hash inside string with comment after",
			input:    `print("Debug: #1") # This should be removed`,
			expected: `print("Debug: #1")`,
		},
		{
			name:     "Triple quoted string",
			input:    `"""This is a docstring # with hash"""`,
			expected: `"""This is a docstring # with hash"""`,
		},
		{
			name:     "Empty line",
			input:    ``,
			expected: ``,
		},
		{
			name:     "Only comment",
			input:    `# This is only a comment`,
			expected: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rc := &RemoveComments{}
			result, keep := rc.removeCommentsFromLine(tc.input)
			
			// If expected is empty, we expect the line to be removed
			if tc.expected == "" && !keep {
				return // Test passed - line was correctly removed
			}
			
			// If expected is not empty, we expect the line to be kept
			if !keep {
				t.Errorf("removeCommentsFromLine(%q) returned keep=false but expected %q", 
					tc.input, tc.expected)
				return
			}
			
			if result != tc.expected {
				t.Errorf("removeCommentsFromLine(%q) = %q; expected %q", 
					tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsDocstring(t *testing.T) {
	testCases := []struct {
		name      string
		line      string
		prevLines []string
		expected  bool
	}{
		{
			name:      "Module docstring at start",
			line:      `"""This is a module docstring"""`,
			prevLines: []string{},
			expected:  true,
		},
		{
			name:      "Function docstring",
			line:      `    """This is a function docstring"""`,
			prevLines: []string{"def add(a, b):"},
			expected:  true,
		},
		{
			name:      "Class docstring",
			line:      `    """This is a class docstring"""`,
			prevLines: []string{"class MyClass:"},
			expected:  true,
		},
		{
			name:      "Not a docstring - regular string",
			line:      `x = """This is not a docstring"""`,
			prevLines: []string{"y = 5"},
			expected:  false,
		},
		{
			name:      "Not a docstring - no triple quotes",
			line:      `print("Hello world")`,
			prevLines: []string{"def test():"},
			expected:  false,
		},
		{
			name:      "Docstring with single quotes",
			line:      `    '''Function docstring'''`,
			prevLines: []string{"def multiply(x, y):"},
			expected:  true,
		},
		{
			name:      "Function docstring with empty lines",
			line:      `    """Calculate sum"""`,
			prevLines: []string{"def calculate():", "", ""},
			expected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rc := &RemoveComments{}
			result := rc.isDocstring(tc.line, tc.prevLines)
			if result != tc.expected {
				t.Errorf("isDocstring(%q, %v) = %v; expected %v",
					tc.line, tc.prevLines, result, tc.expected)
			}
		})
	}
}

func TestRemoveCommentsAndDocstrings(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Simple function with docstring and comments",
			input: `def add(a, b):
    """Add two numbers together."""
    # This is a comment
    return a + b`,
			expected: `def add(a, b):
    return a + b`,
		},
		{
			name: "Module docstring",
			input: `"""This is a module docstring."""
import math

def test():
    pass`,
			expected: `import math

def test():
    pass`,
		},
		{
			name: "Mixed quotes and regular strings",
			input: `def process():
    """Process data."""
    message = "Processing # of items"
    return message`,
			expected: `def process():
    message = "Processing # of items"
    return message`,
		},
		{
			name: "Multi-line docstring",
			input: `def calculate():
    """
    This is a multi-line
    docstring with details.
    """
    result = 42  # Important number
    return result`,
			expected: `def calculate():
    result = 42
    return result`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rc := &RemoveComments{}
			result := rc.removeCommentsAndDocstrings(tc.input)
			if result != tc.expected {
				t.Errorf("removeCommentsAndDocstrings() failed for %s\nGot:\n%q\nExpected:\n%q",
					tc.name, result, tc.expected)
			}
		})
	}
}
