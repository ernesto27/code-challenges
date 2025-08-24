package main

import (
	"testing"
	"time"
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

func TestCPUCoresStats(t *testing.T) {
	cpuStats, err := NewCPUStats()
	if err != nil {
		t.Fatalf("Failed to create CPUStats: %v", err)
	}

	time.Sleep(5 * time.Second) // Wait a second to get different readings
	data, err := cpuStats.CalculateUsagePerCore()
	if err != nil {
		t.Fatalf("Failed to calculate CPU usage per core: %v", err)
	}

	if len(data) == 0 {
		t.Error("No CPU core usage data returned")
	}
}
