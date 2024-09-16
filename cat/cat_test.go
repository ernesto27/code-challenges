package main

import (
	"os"
	"testing"
)

func TestCat_WithFiles(t *testing.T) {
	content := "Hello\nWorld\n"
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cat, err := NewCat().WithFiles(tmpfile.Name())
	if err != nil {
		t.Fatalf("WithFiles failed: %v", err)
	}
	defer cat.Close()

	result, err := cat.GetContentFile(false, false)
	if err != nil {
		t.Fatalf("GetContentFile failed: %v", err)
	}

	if result != content {
		t.Errorf("Expected content %q, got %q", content, result)
	}
}

func TestCat_MultipleFiles(t *testing.T) {
	content1 := "File1\nContent\n"
	content2 := "File2\nContent\n"
	tmpfile1, err := os.CreateTemp("", "example1")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile1.Name())

	tmpfile2, err := os.CreateTemp("", "example2")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile2.Name())

	if _, err := tmpfile1.Write([]byte(content1)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile1.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile2.Write([]byte(content2)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile2.Close(); err != nil {
		t.Fatal(err)
	}

	// Test WithFiles with multiple files
	cat, err := NewCat().WithFiles(tmpfile1.Name(), tmpfile2.Name())
	if err != nil {
		t.Fatalf("WithFiles failed: %v", err)
	}
	defer cat.Close()

	result, err := cat.GetContentFile(false, false)
	if err != nil {
		t.Fatalf("GetContentFile failed: %v", err)
	}

	expected := content1 + content2
	if result != expected {
		t.Errorf("Expected content %q, got %q", expected, result)
	}
}
