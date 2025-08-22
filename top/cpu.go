package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CPUStats struct {
	user    uint64
	nice    uint64
	system  uint64
	idle    uint64
	iowait  uint64
	irq     uint64
	softirq uint64
	steal   uint64
}

type CPUUsageBreakdown struct {
	userPercent   float64
	systemPercent float64
	idlePercent   float64
}

// Read reads CPU statistics from /proc/stat
func (c *CPUStats) Read() error {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return fmt.Errorf("failed to read CPU line")
	}

	line := scanner.Text()
	fields := strings.Fields(line)
	if len(fields) < 8 || fields[0] != "cpu" {
		return fmt.Errorf("invalid CPU line format")
	}

	values := []*uint64{&c.user, &c.nice, &c.system, &c.idle,
		&c.iowait, &c.irq, &c.softirq, &c.steal}

	for i, val := range values {
		parsed, err := strconv.ParseUint(fields[i+1], 10, 64)
		if err != nil {
			return err
		}
		*val = parsed
	}

	return nil
}

// CalculateUsage calculates CPU usage breakdown compared to previous stats
func (c *CPUStats) CalculateUsage(prev *CPUStats) *CPUUsageBreakdown {
	if prev == nil {
		return &CPUUsageBreakdown{0.0, 0.0, 0.0}
	}

	prevTotal := prev.user + prev.nice + prev.system + prev.idle + prev.iowait + prev.irq + prev.softirq + prev.steal
	currTotal := c.user + c.nice + c.system + c.idle + c.iowait + c.irq + c.softirq + c.steal

	totalDiff := currTotal - prevTotal
	if totalDiff == 0 {
		return &CPUUsageBreakdown{0.0, 0.0, 0.0}
	}

	userDiff := (c.user + c.nice) - (prev.user + prev.nice)
	systemDiff := (c.system + c.irq + c.softirq) - (prev.system + prev.irq + prev.softirq)
	idleDiff := c.idle - prev.idle

	return &CPUUsageBreakdown{
		userPercent:   (float64(userDiff) / float64(totalDiff)) * 100.0,
		systemPercent: (float64(systemDiff) / float64(totalDiff)) * 100.0,
		idlePercent:   (float64(idleDiff) / float64(totalDiff)) * 100.0,
	}
}

// getProcessTimes reads /proc/[pid]/stat to get the process's user and system time.
func (c *CPUStats) GetProcessTimes(pid int) (uint64, error) {
	path := fmt.Sprintf("/proc/%d/stat", pid)
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 15 {
		return 0, fmt.Errorf("malformed stat file for pid %d", pid)
	}

	// Field 14 (utime) and 15 (stime). Indices are 13 and 14.
	utime, err := strconv.ParseUint(fields[13], 10, 64)
	if err != nil {
		return 0, err
	}
	stime, err := strconv.ParseUint(fields[14], 10, 64)
	if err != nil {
		return 0, err
	}

	// Total process time is the sum of time spent in user mode and kernel mode.
	return utime + stime, nil
}

// NewCPUStats creates and reads initial CPU statistics
func NewCPUStats() (*CPUStats, error) {
	stats := &CPUStats{}
	err := stats.Read()
	return stats, err
}
