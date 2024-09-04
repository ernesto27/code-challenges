package main

import (
	"os"
	"sort"
	"testing"
)

func TestZipCrackerIsValid(t *testing.T) {
	validFile, err := os.CreateTemp("", "valid_zip_*.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(validFile.Name())

	validSignature := []byte{0x50, 0x4B, 0x03, 0x04}
	_, err = validFile.Write(validSignature)
	if err != nil {
		t.Fatal(err)
	}
	validFile.Close()

	invalidFile, err := os.CreateTemp("", "invalid_zip_*.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(invalidFile.Name())

	invalidSignature := []byte{0x00, 0x00, 0x00, 0x00}
	_, err = invalidFile.Write(invalidSignature)
	if err != nil {
		t.Fatal(err)
	}
	invalidFile.Close()

	tests := []struct {
		name     string
		fileName string
		want     bool
	}{
		{"Valid ZIP", validFile.Name(), true},
		{"Invalid ZIP", invalidFile.Name(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zc, err := NewZipCracker(tt.fileName, "passwords.txt", 0, 0)
			if err != nil {
				t.Fatal(err)
			}
			defer zc.fileZip.Close()

			if got := zc.isValid(); got != tt.want {
				t.Errorf("ZipCracker.isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipCrackerGenerateCombinations(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     []string
	}{
		{
			name:     "Single character",
			password: "a",
			want:     []string{"a"},
		},
		{
			name:     "Two characters",
			password: "ab",
			want:     []string{"ab", "ba"},
		},
		{
			name:     "Three characters",
			password: "abc",
			want:     []string{"abc", "acb", "bac", "bca", "cab", "cba"},
		},
		{
			name:     "Four characters",
			password: "abcd",
			want:     []string{"abcd", "abdc", "acbd", "acdb", "adbc", "adcb", "bacd", "badc", "bcad", "bcda", "bdac", "bdca", "cabd", "cadb", "cbad", "cbda", "cdab", "cdba", "dabc", "dacb", "dbac", "dbca", "dcab", "dcba"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zc := &ZipCracker{}
			got := zc.generateCombinations(tt.password)

			if len(got) != len(tt.want) {
				t.Errorf("ZipCracker.generateCombinations() returned %d combinations, want %d", len(got), len(tt.want))
			}

			sort.Strings(got)
			sort.Strings(tt.want)

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ZipCracker.generateCombinations()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
