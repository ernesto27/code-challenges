package main

import (
	"os"
	"testing"
)

// TODO UPDATE TESTS

func TestCat_WithFiles(t *testing.T) {
	content := `Hello
World
`
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
	content1 := `File1
Content
`
	content2 := `File2
Content
`
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

// New test cases
func TestCat_WithNumbering(t *testing.T) {
	content := `Line 1
Line 2
Line 3
`
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

	result, err := cat.GetContentFile(true, false)
	if err != nil {
		t.Fatalf("GetContentFile failed: %v", err)
	}

	expected := `     1 Line 1
     2 Line 2
     3 Line 3
`
	if result != expected {
		t.Errorf("Expected content %q, got %q", expected, result)
	}
}

func TestCat_WithSqueezeBlank(t *testing.T) {
	content := `Line 1

Line 2

Line 3
`
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

	result, err := cat.GetContentFile(false, true)
	if err != nil {
		t.Fatalf("GetContentFile failed: %v", err)
	}

	expected := `     1 Line 1

     2 Line 2

     3 Line 3
`
	if result != expected {
		t.Errorf("Expected content %q, got %q", expected, result)
	}
}

func TestCat_WithNumberingWithoutIgnoringEmptyLines(t *testing.T) {
	content := `Line 1

Line 3

Line 5
`
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

	result, err := cat.GetContentFile(true, false)
	if err != nil {
		t.Fatalf("GetContentFile failed: %v", err)
	}

	expected := `     1 Line 1
     2 
     3 Line 3
     4 
     5 Line 5
`
	if result != expected {
		t.Errorf("Expected content %q, got %q", expected, result)
	}
}

func TestCat_EmptyFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "empty")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	cat, err := NewCat().WithFiles(tmpfile.Name())
	if err != nil {
		t.Fatalf("WithFiles failed: %v", err)
	}
	defer cat.Close()

	result, err := cat.GetContentFile(false, false)
	if err != nil {
		t.Fatalf("GetContentFile failed: %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty content, got %q", result)
	}
}

func TestCat_NonExistentFile(t *testing.T) {
	_, err := NewCat().WithFiles("non_existent_file.txt")
	if err == nil {
		t.Error("Expected an error for non-existent file, got nil")
	}
}
