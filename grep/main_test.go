package main

import (
	"testing"
)

func TestRunEmptyFilter(t *testing.T) {
	grep := Grep{}

	content := "Hello\nWorld\n"
	filter := ""

	result, err := grep.Run(filter, content, false)
	if err != nil {
		t.Error("Error should be nil")
	}

	if result != content {
		t.Error("Result should be equal to content")
	}
}

func TestRunWithFilter(t *testing.T) {
	const redColor = "\033[31m"
	const resetColor = "\033[0m"

	grep := Grep{}

	content := "Metallica\nNirvana\nMegadeth"
	filter := "M"

	result, err := grep.Run(filter, content, false)
	if err != nil {
		t.Error("Error should be nil")
	}

	expected := redColor + "M" + resetColor + "etallica\n" + redColor + "M" + resetColor + "egadeth"

	if result != expected {
		t.Errorf("Result should be equal to %s, got %s", expected, result)
	}
}
