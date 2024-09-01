package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	text := `line1
line2
line2
line3
line4`

	expected := `line1
line2
line3
line4`

	uniq := Uniq{}
	resp := uniq.Run(text)
	if resp != expected {
		t.Errorf("main() = %s; expected %s", resp, expected)
	}
}
func TestSave(t *testing.T) {
	text := `line1
line2
line2
line3
line4`
	expected := `line1
line2
line3
line4`
	uniq := Uniq{}
	uniq.Run(text)

	err := uniq.Save("test_output.txt")
	if err != nil {
		t.Errorf("Save() returned an error: %v", err)
	}

	fileContent, err := os.ReadFile("test_output.txt")
	if err != nil {
		t.Errorf("Failed to read test_output.txt: %v", err)
	}

	actual := string(fileContent)
	if actual != expected {
		t.Errorf("Save() wrote incorrect content to file. Got:\n%s\nExpected:\n%s", actual, expected)
	}

	err = os.Remove("test_output.txt")
	if err != nil {
		t.Errorf("Failed to remove test_output.txt: %v", err)
	}
}
func TestCount(t *testing.T) {
	text := `line1
line2
line2
line3
line1
line4
line2
line4
line5`
	expected := `2 line1
3 line2
1 line3
2 line4
1 line5`
	uniq := Uniq{}
	actual := uniq.Count(text)
	if actual != expected {
		t.Errorf("Count() = %s; expected %s", actual, expected)
	}
}

func TestGetDuplicated(t *testing.T) {
	text := `line1
line2
line2
line3
line5
line1
line2
line3`

	expected := `line2
line1
line3`

	uniq := Uniq{}
	actual := uniq.getDuplicated(text)
	if actual != expected {
		t.Errorf("getDuplicated() = %s; expected %s", actual, expected)
	}
}
