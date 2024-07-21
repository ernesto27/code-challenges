package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestDoReplace(t *testing.T) {
	testCases := []struct {
		cmd           string
		text          string
		expected      string
		expectedError bool
	}{
		{
			cmd:           "s/this/that/g",
			text:          "this is a test. this is only a test.",
			expected:      "that is a test. that is only a test.",
			expectedError: false,
		},
		{
			cmd:           "s/\"//g",
			text:          "\"Your heart is the size of an ocean. Go find yourself in is hidden depths size.\"",
			expected:      "Your heart is the size of an ocean. Go find yourself in is hidden depths size.",
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		sed, err := newSed(tc.cmd, tc.text)
		if err != nil {
			if !tc.expectedError {
				t.Fatalf("newSed(%s, %s) returned error: %v", tc.cmd, tc.text, err)
			}
			continue
		}

		actual := sed.doReplace()
		if actual != tc.expected {
			t.Errorf("doReplace() = %s; expected %s", actual, tc.expected)
		}
	}
}

func TestRangePrint(t *testing.T) {
	text := `1	"Your heart is the size of an ocean. Go find yourself in its hidden depths."
	2	"The Bay of Bengal is hit frequently by cyclones. The months of November and May, in particular, are dangerous in this regard."
	3	"Thinking is the capital, Enterprise is the way, Hard Work is the solution."
	4	"If You Can'T Make It Good, At Least Make It Look Good."
	5	"Heart be brave. If you cannot be brave, just go. Love's glory is not a small thing."
	6	"It is bad for a young man to sin; but it is worse for an old man to sin."
	7	"If You Are Out To Describe The Truth, Leave Elegance To The Tailor."
	8	"O man you are busy working for the world, and the world is busy trying to turn you out."
	9	"While children are struggling to be unique, the world around them is trying all means to make them look like everybody else."
	10	"These Capitalists Generally Act Harmoniously And In Concert, To Fleece The People."`

	sed, err := newSed("2,4p", text)
	if err != nil {
		t.Fatalf("newSed(2,4p, some text) returned error: %v", err)
	}

	resp := sed.rangeLines()

	lines := strings.Split(resp, "\n")

	if len(lines) != 3 {
		t.Fatalf("rangePrintLines() returned %d lines; expected 4", len(lines))
	}
}

func TestGetLinesByPattern(t *testing.T) {

	text := `1	"Your heart is the size of an ocean. Go find yourself in its hidden depths."
	2	"The Bay of Bengal is hit frequently by cyclones. The months of November and May, in particular, are dangerous in this regard."
	3	"Thinking is the capital, Enterprise is the way, Hard Work is the solution."
	4	"If You Can'T Make It Good, At Least Make It Look Good."
	5	"Heart be brave. If you cannot be brave, just go. Love's glory is not a small thing."
	6	"It is bad for a young"`

	sed, err := newSed("/heart/p", text)
	if err != nil {
		t.Fatalf("newSed(/Heart/p, some text) returned error: %v", err)
	}

	resp := sed.getLinesByPattern()

	lines := strings.Split(resp, "\n")
	if len(lines) != 1 {
		t.Fatalf("getLinesByPattern() returned %d lines; expected 1", len(lines))
	}
}

func TestDoubleSpacing(t *testing.T) {
	text := `Your heart is the size of an ocean. Go find yourself in its hidden depths."
The Bay of Bengal is hit frequently by cyclones. The months of November and May, in particular, are dangerous in this regard."
Thinking is the capital, Enterprise is the way, Hard Work is the solution."`

	sed, err := newSed("G", text)
	if err != nil {
		t.Fatalf("newSed(G, some text) returned error: %v", err)
	}

	resp := sed.doubleSpacing()

	fmt.Println(resp)
	lines := strings.Split(resp, "\n")

	if len(lines) != 5 {
		t.Fatalf("doubleSpacing() returned %d lines; expected 5", len(lines))
	}
}

func TestRemoveTrailingspaces(t *testing.T) {
	text := `	Your heart is the size of an ocean. Go find yourself in its hidden depths size.		`

	sed, err := newSed("G", text)
	if err != nil {
		t.Fatalf("newSed(G, some text) returned error: %v", err)
	}

	resp := sed.removeTrailingSpaces()

	if resp != "Your heart is the size of an ocean. Go find yourself in its hidden depths size." {
		t.Fatalf("removeTrailingSpaces() returned %s; expected 'Your heart is the size of an ocean. Go find yourself in its hidden depths size.'", resp)
	}
}
