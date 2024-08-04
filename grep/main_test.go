package main

import (
	"testing"
)

type GrepTest struct {
	name            string
	filter          string
	content         string
	exclude         bool
	caseInsensitive bool
	expected        string
}

func TestGrep_Run(t *testing.T) {
	tests := []GrepTest{
		{
			name:            "Empty Filter",
			filter:          "",
			content:         "Hello\nWorld\n",
			exclude:         false,
			caseInsensitive: false,
			expected:        "Hello\nWorld\n",
		},
		{
			name:            "Filter with 'M'",
			filter:          "M",
			content:         "Metallica\nNirvana\nMegadeth",
			exclude:         false,
			caseInsensitive: false,
			expected:        redColor + "M" + resetColor + "etallica\n" + redColor + "M" + resetColor + "egadeth",
		},
		{
			name:            "Exclude Filter with 'M'",
			filter:          "M",
			content:         "Metallica\nNirvana\nMegadeth",
			exclude:         true,
			caseInsensitive: false,
			expected:        "Nirvana",
		},
		{
			name:            "Filter with \\d to match numbers",
			filter:          "\\d+",
			content:         "Metallica\nNirvana\nMegadeth\n1",
			exclude:         false,
			caseInsensitive: false,
			expected:        redColor + "1" + resetColor,
		},
		{
			name:            "Filter with \\w to match words",
			filter:          "\\w+",
			content:         "!\nmystring\n@",
			exclude:         false,
			caseInsensitive: false,
			expected:        redColor + "mystring" + resetColor,
		},
		{
			name:            "Case Insensitive Filter with 'm'",
			filter:          "m",
			content:         "Metallica\nNirvana\nMegadeth",
			exclude:         false,
			caseInsensitive: true,
			expected:        redColor + "M" + resetColor + "etallica\n" + redColor + "M" + resetColor + "egadeth",
		},
	}

	grep := Grep{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := grep.Run(test.filter, test.content, test.exclude, test.caseInsensitive)
			if err != nil {
				t.Error("Error should be nil")
			}

			if result != test.expected {
				t.Errorf("Result should be equal to %s, got %s", test.expected, result)
			}
		})
	}
}

func TestRecursive(t *testing.T) {
	grep := Grep{}

	initialPath := "./unit-test"
	filter := "Nirvana"

	result, err := grep.RunRecursive(filter, initialPath)
	if err != nil {
		t.Error("Error should be nil")
	}

	expected := "unit-test/custom.txt:" + redColor + "Nirvana" + resetColor + " text\nunit-test/custom.txt:" + redColor + "Nirvana" + resetColor + " text2"

	if result != expected {
		t.Errorf("Result should be equal to %s, got %s", expected, result)
	}
}
