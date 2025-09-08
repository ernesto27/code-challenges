package main

import (
	"testing"
)

func TestJavaScriptCommentRemoval(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Single line comments",
			input: `const x = 5; // This is a comment
const y = 10; // Another comment
const z = x + y;`,
			expected: `const x = 5;
const y = 10;
const z = x + y;`,
		},
		{
			name: "Multi-line block comments",
			input: `/* This is a 
   multi-line comment */
const result = calculate();`,
			expected: `const result = calculate();`,
		},
		{
			name: "JSDoc comments",
			input: `/**
 * This is a JSDoc comment
 * @param {number} a - First parameter
 */
function add(a, b) {
    return a + b;
}`,
			expected: `function add(a, b) {
    return a + b;
}`,
		},
		{
			name: "Comments inside strings should be preserved",
			input: `const message = "This // is not a comment";
const url = "https://example.com /* not a comment */";
const code = 'console.log("// still not a comment");';`,
			expected: `const message = "This // is not a comment";
const url = "https://example.com /* not a comment */";
const code = 'console.log("// still not a comment");';`,
		},
		{
			name: "Mixed comment types",
			input: `// Single line comment
/* Block comment */ const x = 5;
function test() {
    return x; // Inline comment
}`,
			expected: `const x = 5;
function test() {
    return x;
}`,
		},
		// {
		// 	name:     "Template literals with comments",
		// 	input:    "const template = `This is a template // not a comment\nwith multiple lines /* also not a comment */`;\n// This is a real comment",
		// 	expected: "const template = `This is a template // not a comment\nwith multiple lines /* also not a comment */`;",
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			rc := &RemoveCommentsC{
				nameFile:        "test.js",
				originalContent: tc.input,
			}
			result := rc.Remove()
			if result != tc.expected {
				t.Errorf("Remove() failed for %s\nGot:\n%q\nExpected:\n%q",
					tc.name, result, tc.expected)
			}
		})
	}
}

func TestJavaScriptIsInsideString(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		pos      int
		expected bool
	}{
		{
			name:     "Double slash outside string",
			line:     `const x = 5; // comment`,
			pos:      13,
			expected: false,
		},
		{
			name:     "Double slash inside double quotes",
			line:     `const url = "http://example.com";`,
			pos:      19,
			expected: true,
		},
		{
			name:     "Double slash inside single quotes",
			line:     `const url = 'http://example.com';`,
			pos:      19,
			expected: true,
		},
		{
			name:     "Double slash inside template literal",
			line:     "const url = `http://example.com`;",
			pos:      19,
			expected: true,
		},
		{
			name:     "Block comment start outside string",
			line:     `const x = 5; /* comment */`,
			pos:      13,
			expected: false,
		},
		{
			name:     "Block comment start inside string",
			line:     `const text = "/* not a comment */";`,
			pos:      15,
			expected: true,
		},
		{
			name:     "Escaped quotes in string",
			line:     `const text = "He said \"Hi\" // test";`,
			pos:      25,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rc := &RemoveCommentsC{
				nameFile:        "test.js",
				originalContent: "",
			}
			result := rc.isInsideString(tc.line, tc.pos)
			if result != tc.expected {
				t.Errorf("isInsideString(%q, %d) = %v; expected %v",
					tc.line, tc.pos, result, tc.expected)
			}
		})
	}
}
