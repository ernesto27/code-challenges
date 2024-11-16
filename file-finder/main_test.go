package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestFind(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filetest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFiles := []struct {
		name string
		size int64
	}{
		{"small.txt", 500},
		{"large.txt", 1500},
	}

	for _, tf := range testFiles {
		path := filepath.Join(tmpDir, tf.name)
		data := make([]byte, tf.size)
		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name       string
		minSize    int64
		wantFiles  int
		wantErrNil bool
	}{
		{
			name:       "find large files",
			minSize:    1000,
			wantFiles:  1,
			wantErrNil: true,
		},
		{
			name:       "find all files",
			minSize:    0,
			wantFiles:  2,
			wantErrNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := NewFileFinder(tmpDir, tt.minSize)
			err := ff.Find()

			if (err == nil) != tt.wantErrNil {
				t.Errorf("Find() error = %v, wantErrNil %v", err, tt.wantErrNil)
			}

			if len(ff.files) != tt.wantFiles {
				t.Errorf("Find() found %v files, want %v", len(ff.files), tt.wantFiles)
			}
		})
	}
}

func TestGenerateMD5FromFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "filetest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	data := []byte("Hello, World!")
	if _, err := tmpFile.Write(data); err != nil {
		t.Fatal(err)
	}

	fileFinder := NewFileFinder("", 0)

	hash, err := fileFinder.generateMD5(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	want := "65a8e27d8879283831b664bd8b7f0ad4"
	if hash != want {
		t.Errorf("generateMD5FromFile() = %v, want %v", hash, want)
	}
}

func TestCompareByteByByte(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "comparetest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testCases := []struct {
		name     string
		content1 []byte
		content2 []byte
		want     bool
		wantErr  bool
	}{
		{
			name:     "identical files",
			content1: []byte("hello world"),
			content2: []byte("hello world"),
			want:     true,
			wantErr:  false,
		},
		{
			name:     "different files",
			content1: []byte("hello world"),
			content2: []byte("hello world!"),
			want:     false,
			wantErr:  false,
		},
		{
			name:     "empty files",
			content1: []byte(""),
			content2: []byte(""),
			want:     true,
			wantErr:  false,
		},
		{
			name:     "different sizes",
			content1: bytes.Repeat([]byte("a"), 100000),
			content2: bytes.Repeat([]byte("a"), 90000),
			want:     false,
			wantErr:  false,
		},
	}

	ff := NewFileFinder(".", 0)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file1 := filepath.Join(tmpDir, "file1_"+tc.name)
			file2 := filepath.Join(tmpDir, "file2_"+tc.name)

			if err := os.WriteFile(file1, tc.content1, 0644); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(file2, tc.content2, 0644); err != nil {
				t.Fatal(err)
			}

			got, err := ff.compareByteByByte(file1, file2)
			if (err != nil) != tc.wantErr {
				t.Errorf("compareByteByByte() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("compareByteByByte() = %v, want %v", got, tc.want)
			}
		})
	}
}
