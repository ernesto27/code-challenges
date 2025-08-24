package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CPUCoreStats struct {
	user    uint64
	nice    uint64
	system  uint64
	idle    uint64
	iowait  uint64
	irq     uint64
	softirq uint64
	steal   uint64
}

type CPUStats struct {
	CPUCoreStats
	cores []CPUCoreStats
}

type CPUUsageBreakdown struct {
	userPercent   float64
	systemPercent float64
	idlePercent   float64
}

// parseCPUFields parses CPU statistics fields from /proc/stat line
func parseCPUFields(fields []string, values []*uint64) error {
	if len(fields) < len(values)+1 {
		return fmt.Errorf("insufficient fields in CPU line")
	}

	for i, val := range values {
		parsed, err := strconv.ParseUint(fields[i+1], 10, 64)
		if err != nil {
			return err
		}
		*val = parsed
	}
	return nil
}

// Read reads CPU statistics from /proc/stat
func (c *CPUStats) Read() error {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read individual CPU core stats
	c.cores = []CPUCoreStats{}
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if strings.Contains(line, "cpu") {

			if len(fields) < 8 || !strings.HasPrefix(fields[0], "cpu") || fields[0] == "cpu" {
				continue
			}

			core := CPUCoreStats{}
			coreValues := []*uint64{&core.user, &core.nice, &core.system, &core.idle,
				&core.iowait, &core.irq, &core.softirq, &core.steal}

			if err := parseCPUFields(fields, coreValues); err != nil {
				return err
			}

			c.cores = append(c.cores, core)
		}
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

func (c *CPUStats) CalculateUsagePerCore() (map[int]float64, error) {
	resp := map[int]float64{}

	newCpuStats := &CPUStats{}
	err := newCpuStats.Read()
	if err != nil {
		return resp, err
	}

	for i := range c.cores {
		prevTotal := c.cores[i].user + c.cores[i].nice + c.cores[i].system + c.cores[i].idle +
			c.cores[i].iowait + c.cores[i].irq + c.cores[i].softirq + c.cores[i].steal

		currTotal := newCpuStats.cores[i].user + newCpuStats.cores[i].nice + newCpuStats.cores[i].system + newCpuStats.cores[i].idle +
			newCpuStats.cores[i].iowait + newCpuStats.cores[i].irq + newCpuStats.cores[i].softirq + newCpuStats.cores[i].steal

		prevIdle := c.cores[i].idle + c.cores[i].iowait
		currIdle := newCpuStats.cores[i].idle + newCpuStats.cores[i].iowait

		deltaTotal := float64(currTotal - prevTotal)
		deltaIdle := float64(currIdle - prevIdle)

		if deltaTotal == 0 {
			resp[i] = 0.0
		} else {
			usage := (deltaTotal - deltaIdle) / deltaTotal * 100.0
			resp[i] = usage
		}
	}

	// Update current stats with new readings
	*c = *newCpuStats

	return resp, nil
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
