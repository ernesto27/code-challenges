package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestListFiles(t *testing.T) {
	tarFile, err := NewTar("files.tar", false)
	if err != nil {
		t.Fatal(err)
	}
	defer tarFile.Close()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = tarFile.ListFiles()
	if err != nil {
		t.Fatal(err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Verify the output
	expectedOutput := "file1.txt\nfile2.txt\nfile3.txt\n"
	if buf.String() != expectedOutput {
		t.Errorf("expected %q, got %q", expectedOutput, buf.String())
	}
}

func TestCreateNewTar(t *testing.T) {
	tarFile, err := NewTar("newfile.tar", true)
	if err != nil {
		t.Fatal(err)
	}
	defer tarFile.Close()
	defer os.Remove("newfile.tar")

	err = tarFile.CreateTar([]string{"file1.txt", "file2.txt", "file3.txt"})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat("newfile.tar"); os.IsNotExist(err) {
		t.Errorf("expected newfile.tar to exist, but it does not")
	}
}

func TestExtractTar(t *testing.T) {
	tarFile, err := NewTar("files.tar", false)
	if err != nil {
		t.Fatal(err)
	}
	defer tarFile.Close()

	pathToExtract := "test"
	err = tarFile.ExtractFiles("test")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(pathToExtract + "/file1.txt")
	if os.IsNotExist(err) {
		t.Errorf(err.Error())
	}

	_, err = os.Stat(pathToExtract + "/file1.txt")
	if os.IsNotExist(err) {
		t.Errorf(err.Error())
	}

	_, err = os.Stat(pathToExtract + "/file1.txt")
	if os.IsNotExist(err) {
		t.Errorf(err.Error())
	}
}
