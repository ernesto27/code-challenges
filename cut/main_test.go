package main

import (
	"os"
	"reflect"
	"testing"
)

func TestCutRead(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		index     string
		delimiter string
		expected  map[int][]string
		expectErr bool
	}{
		{
			name:      "Basic CSV",
			data:      "a,b,c\n1,2,3\n4,5,6",
			index:     "1,3",
			delimiter: ",",
			expected: map[int][]string{
				0: {"a", "1", "4"},
				2: {"c", "3", "6"},
			},
			expectErr: false,
		},
		{
			name:      "Tab-separated",
			data:      "a\tb\tc\n1\t2\t3\n4\t5\t6",
			index:     "2",
			delimiter: "\t",
			expected: map[int][]string{
				1: {"b", "2", "5"},
			},
			expectErr: false,
		},
		{
			name:      "Single column",
			data:      "a,b,c\n1,2,3\n4,5,6",
			index:     "2",
			delimiter: ",",
			expected: map[int][]string{
				1: {"b", "2", "5"},
			},
			expectErr: false,
		},
		{
			name:      "Out of range index",
			data:      "a,b,c\n1,2,3\n4,5,6",
			index:     "4",
			delimiter: ",",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "Empty file",
			data:      "",
			index:     "1",
			delimiter: ",",
			expected:  map[int][]string{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			_, err = tmpfile.Write([]byte(tt.data))
			if err != nil {
				t.Fatal(err)
			}
			tmpfile.Close()

			cut := NewCut(tmpfile.Name())
			if cut == nil {
				t.Fatal("NewCut returned nil")
			}

			err = cut.Read(tt.index, tt.delimiter)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(cut.recordList, tt.expected) {
					t.Errorf("Expected recordList %v, got %v", tt.expected, cut.recordList)
				}
			}
		})
	}
}

func TestCutReadInvalidIndex(t *testing.T) {
	cut := NewCutWithStdin()
	err := cut.Read("0", ",")
	if err == nil || err.Error() != "not valid index" {
		t.Errorf("Expected 'not valid index' error, got %v", err)
	}
}
