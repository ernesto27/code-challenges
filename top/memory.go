package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MemoryInfo struct {
	totalKB     uint64
	availableKB uint64
	usageInGB   string
}

// GetUsageInGB returns memory usage in format like "12GB/32GB"
func (m *MemoryInfo) GetUsageInGB() (string, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	fieldsFound := 0
	targetFields := 2

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		valueStr := fields[1]
		value, err := strconv.ParseUint(valueStr, 10, 64)
		if err != nil {
			continue
		}

		switch key {
		case "MemTotal":
			m.totalKB = value
			fieldsFound++
		case "MemAvailable":
			m.availableKB = value
			fieldsFound++
		}

		// Stop early if we have all target values
		if fieldsFound >= targetFields {
			break
		}
	}

	usedKB := m.totalKB - m.availableKB
	usedGB := float64(usedKB) / 1024 / 1024
	totalGB := float64(m.totalKB) / 1024 / 1024

	usageInGB := fmt.Sprintf("%.1fGB/%.1fGB", usedGB, totalGB)

	m.usageInGB = usageInGB

	return usageInGB, nil
}

func (m *MemoryInfo) GetPercentageUse() float64 {
	used := m.totalKB - m.availableKB
	percentage := 0.0
	if m.totalKB > 0 {
		percentage = float64(used) / float64(m.totalKB)
	}

	return percentage

}

// NewMemoryInfo creates and reads initial memory information
func NewMemoryInfo() *MemoryInfo {
	info := &MemoryInfo{}
	return info
}
