package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MemoryInfo struct {
	totalKB       uint64
	availableKB   uint64
	freeKB        uint64
	buffersKB     uint64
	cachedKB      uint64
	unevictableKB uint64
	activeKB      uint64
	swapCachedKB  uint64
}

// Read reads memory information from /proc/meminfo
func (m *MemoryInfo) Read() error {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	fieldsFound := 0
	targetFields := 8 // Number of fields we want to collect

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
		case "MemFree":
			m.freeKB = value
			fieldsFound++
		case "Buffers":
			m.buffersKB = value
			fieldsFound++
		case "Cached":
			m.cachedKB = value
			fieldsFound++
		case "Unevictable":
			m.unevictableKB = value
			fieldsFound++
		case "Active":
			m.activeKB = value
			fieldsFound++
		case "SwapCached":
			m.swapCachedKB = value
			fieldsFound++
		}

		// Stop early if we have all target values
		if fieldsFound >= targetFields {
			break
		}
	}

	return scanner.Err()
}

// GetUsedMemory calculates used memory in KB
func (m *MemoryInfo) GetUsedMemory() uint64 {
	return m.totalKB - m.availableKB
}

// GetWiredMemory calculates approximate "wired" memory in KB
func (m *MemoryInfo) GetWiredMemory() uint64 {
	return m.unevictableKB + (m.activeKB / 4)
}

// GetCompressorMemory calculates approximate "compressor" memory in KB
func (m *MemoryInfo) GetCompressorMemory() uint64 {
	return m.buffersKB + m.swapCachedKB + m.cachedKB/2
}

// GetUnusedMemory returns free memory in KB
func (m *MemoryInfo) GetUnusedMemory() uint64 {
	return m.freeKB
}

// FormatSize formats memory size from KB to human-readable format
func (m *MemoryInfo) FormatSize(kb uint64) string {
	if kb >= 1024*1024 {
		gb := float64(kb) / (1024 * 1024)
		return fmt.Sprintf("%.0fGB", gb)
	} else if kb >= 1024 {
		mb := float64(kb) / 1024
		return fmt.Sprintf("%.0fMB", mb)
	}
	return fmt.Sprintf("%dKB", kb)
}

// NewMemoryInfo creates and reads initial memory information
func NewMemoryInfo() (*MemoryInfo, error) {
	info := &MemoryInfo{}
	err := info.Read()
	return info, err
}