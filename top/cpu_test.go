package main

import (
	"testing"
)

func TestCPUStats_Read(t *testing.T) {
	stats := &CPUStats{}
	err := stats.Read()
	if err != nil {
		t.Fatalf("Failed to read CPU stats: %v", err)
	}

	// Verify that all fields are non-zero (they should have some values on a running system)
	if stats.user == 0 && stats.nice == 0 && stats.system == 0 && stats.idle == 0 {
		t.Error("All CPU stats are zero, which is unexpected")
	}
}
