package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestHead_Run(t *testing.T) {
	tests := []struct {
		name     string
		content  []string
		lines    int
		bytes    int
		expected string
	}{
		{
			name:     "Read first 2 lines",
			content:  []string{"line1\nline2\nline3\n"},
			lines:    2,
			bytes:    0,
			expected: "line1\nline2\n",
		},
		{
			name:     "Read first 5 bytes",
			content:  []string{"line1\nline2\nline3\n"},
			lines:    0,
			bytes:    5,
			expected: "line1",
		},
		{
			name:     "Read all lines if less than specified",
			content:  []string{"line1\nline2\n"},
			lines:    5,
			bytes:    0,
			expected: "line1\nline2\n",
		},
		{
			name:     "Read all bytes if less than specified",
			content:  []string{"line1\nline2\n"},
			lines:    0,
			bytes:    20,
			expected: "line1\nline2\n",
		},
		{
			name:     "Read first 2 lines from multiple files",
			content:  []string{"file1_line1\nfile1_line2\nfile1_line3\n", "file2_line1\nfile2_line2\nfile2_line3\n"},
			lines:    2,
			bytes:    0,
			expected: "==> /tmp/example1 <==\nfile1_line1\nfile1_line2\n==> /tmp/example2 <==\nfile2_line1\nfile2_line2\n",
		},
		{
			name:     "Read first 5 bytes from multiple files",
			content:  []string{"file1_line1\nfile1_line2\nfile1_line3\n", "file2_line1\nfile2_line2\nfile2_line3\n"},
			lines:    0,
			bytes:    5,
			expected: "file1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tmpfiles []*os.File
			for _, content := range tt.content {
				tmpfile, err := os.CreateTemp("", "example")
				if err != nil {
					t.Fatal(err)
				}
				defer os.Remove(tmpfile.Name())

				if _, err := tmpfile.WriteString(content); err != nil {
					t.Fatal(err)
				}
				if err := tmpfile.Close(); err != nil {
					t.Fatal(err)
				}

				tmpfiles = append(tmpfiles, tmpfile)
			}

			var fileNames []string
			for _, tmpfile := range tmpfiles {
				fileNames = append(fileNames, tmpfile.Name())
			}

			head, err := newHead(fileNames)
			if err != nil {
				t.Fatal(err)
			}
			defer head.Close()

			result, err := head.Run(tt.lines, tt.bytes)
			if err != nil {
				t.Fatal(err)
			}

			for i, fileName := range fileNames {
				placeholder := fmt.Sprintf("/tmp/example%d", i+1)
				result = strings.Replace(result, fileName, placeholder, -1)
			}

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
