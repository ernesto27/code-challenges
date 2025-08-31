package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var processMap = make(map[int]*ProcessInfo)

type ProcessInfo struct {
	PID       int
	PPID      int
	Name      string
	State     string
	CPUUsage  float64
	MemoryKB  uint64
	UTime     uint64
	STime     uint64
	StartTime uint64
	UID       uint32
	Username  string

	// Previous measurement for CPU calculation
	PrevUTime           uint64
	PrevSTime           uint64
	PrevMeasurementTime int64
	HasPrevMeasurement  bool

	ThreadsCount int
}

type ProcessStats struct {
	TotalCPUTime uint64
	SystemTime   uint64
}

// Read populates process information from /proc/[pid]/
func (p *ProcessInfo) Read() error {

	// Read process name from /proc/[pid]/comm
	if err := p.readComm(); err != nil {
		return err
	}

	// Read detailed stats from /proc/[pid]/stat
	if err := p.readStat(); err != nil {
		return err
	}

	// Read memory usage from /proc/[pid]/status
	if err := p.readStatus(); err != nil {
		return err
	}

	// Read thread count from /proc/[pid]/task
	if err := p.ReadTasks(); err != nil {
		return err
	}

	return nil
}

// readComm reads the process command name
func (p *ProcessInfo) readComm() error {
	commPath := filepath.Join("/proc", strconv.Itoa(p.PID), "comm")
	data, err := os.ReadFile(commPath)
	if err != nil {
		return err
	}
	p.Name = strings.TrimSpace(string(data))
	return nil
}

// readStat reads process statistics from /proc/[pid]/stat
func (p *ProcessInfo) readStat() error {
	statPath := filepath.Join("/proc", strconv.Itoa(p.PID), "stat")
	data, err := os.ReadFile(statPath)
	if err != nil {
		return err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 22 {
		return fmt.Errorf("malformed stat file for pid %d", p.PID)
	}

	// Parse PPID (field 4, index 3)
	ppid, err := strconv.Atoi(fields[3])
	if err != nil {
		return err
	}
	p.PPID = ppid

	// Parse state (field 3, index 2)
	p.State = fields[2]

	// Parse CPU times (fields 14-15, indices 13-14)
	utime, err := strconv.ParseUint(fields[13], 10, 64)
	if err != nil {
		return err
	}
	p.UTime = utime

	stime, err := strconv.ParseUint(fields[14], 10, 64)
	if err != nil {
		return err
	}
	p.STime = stime

	// Parse start time (field 22, index 21)
	starttime, err := strconv.ParseUint(fields[21], 10, 64)
	if err != nil {
		return err
	}
	p.StartTime = starttime

	return nil
}

// readStatus reads memory information, UID, and username from /proc/[pid]/status
func (p *ProcessInfo) readStatus() error {
	statusPath := filepath.Join("/proc", strconv.Itoa(p.PID), "status")
	file, err := os.Open(statusPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var foundMemory, foundUID bool
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")

		if key == "VmRSS" && len(fields) >= 2 && !foundMemory {
			value, err := strconv.ParseUint(fields[1], 10, 64)
			if err == nil {
				p.MemoryKB = value
				foundMemory = true
			}
		} else if key == "Uid" && len(fields) >= 2 && !foundUID {
			// Uid field format: "Uid: real_uid effective_uid saved_uid filesystem_uid"
			// We want the effective UID (second field)
			uid, err := strconv.ParseUint(fields[1], 10, 32)
			if err == nil {
				p.UID = uint32(uid)
				// Resolve username from UID
				if u, err := user.LookupId(fields[1]); err == nil {
					p.Username = u.Username
				} else {
					p.Username = fields[1] // fallback to UID string if lookup fails
				}
				foundUID = true
			}
		}

		// Exit early if we found both values we're looking for
		if foundMemory && foundUID {
			break
		}
	}

	return scanner.Err()
}

func (p *ProcessInfo) ReadTasks() error {
	path := "/proc/" + strconv.Itoa(p.PID) + "/task"

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var dirCount int
	for _, entry := range entries {
		if entry.IsDir() {
			dirCount++
		}
	}
	p.ThreadsCount = dirCount
	return nil

}

// UpdateMeasurement stores current CPU times as previous for next calculation
func (p *ProcessInfo) UpdateMeasurement() {
	p.PrevUTime = p.UTime
	p.PrevSTime = p.STime
	p.PrevMeasurementTime = time.Now().UnixNano()
	p.HasPrevMeasurement = true
}

// CalculateCPUUsage calculates CPU usage percentage using internal previous measurements
func (p *ProcessInfo) CalculateCPUUsage() {
	if !p.HasPrevMeasurement {
		p.CPUUsage = 0.0
		return
	}

	currentTime := time.Now().UnixNano()
	timeDeltaNs := currentTime - p.PrevMeasurementTime

	if timeDeltaNs <= 0 {
		p.CPUUsage = 0.0
		return
	}

	// Convert time delta to seconds
	timeDeltaSeconds := float64(timeDeltaNs) / 1e9

	// CPU time is in clock ticks, typically 100 ticks per second (USER_HZ)
	const clockTicksPerSecond = 100.0

	processCPUDelta := (p.UTime + p.STime) - (p.PrevUTime + p.PrevSTime)
	processCPUSeconds := float64(processCPUDelta) / clockTicksPerSecond

	// Calculate percentage: (CPU time used / wall time elapsed) * 100
	p.CPUUsage = (processCPUSeconds / timeDeltaSeconds) * 100.0

	// Cap at reasonable value (shouldn't exceed 100% per core, but can be higher on multi-core)
	if p.CPUUsage > 999.9 {
		p.CPUUsage = 999.9
	}
}

// FormatMemory formats memory size from KB to human-readable format
func (p *ProcessInfo) FormatMemory() string {
	if p.MemoryKB >= 1024*1024 {
		gb := float64(p.MemoryKB) / (1024 * 1024)
		return fmt.Sprintf("%.1fGB", gb)
	} else if p.MemoryKB >= 1024 {
		mb := float64(p.MemoryKB) / 1024
		return fmt.Sprintf("%.0fMB", mb)
	} else if p.MemoryKB > 0 {
		return fmt.Sprintf("%dKB", p.MemoryKB)
	}
	return "0"
}

// GetTotalCPUTime returns the total CPU time used by this process
func (p *ProcessInfo) GetTotalCPUTime() uint64 {
	return p.UTime + p.STime
}

// IsRunning checks if the process state indicates it's actively running
func (p *ProcessInfo) IsRunning() bool {
	return p.State == "R"
}

// IsSleeping checks if the process is sleeping
func (p *ProcessInfo) IsSleeping() bool {
	return p.State == "S" || p.State == "D"
}

func (p *ProcessInfo) GetTime() string {
	sum := p.UTime + p.STime
	duration := time.Duration(sum) * time.Second

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

}

// GetRunningProcesses returns a list of all running processes with CPU calculations
func GetRunningProcesses() ([]*ProcessInfo, error) {
	entries, err := os.ReadDir("/proc/")
	if err != nil {
		return nil, err
	}

	var processes []*ProcessInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		proc, err := NewProcessInfo(pid)
		if err != nil {
			continue
		}

		processes = append(processes, proc)
	}

	// Update CPU calculations
	for _, proc := range processes {
		if prevProc, exists := processMap[proc.PID]; exists {
			// Copy previous measurements
			proc.PrevUTime = prevProc.UTime
			proc.PrevSTime = prevProc.STime
			proc.PrevMeasurementTime = prevProc.PrevMeasurementTime
			proc.HasPrevMeasurement = prevProc.HasPrevMeasurement
			proc.CalculateCPUUsage()
		}
		proc.UpdateMeasurement()
		processCopy := *proc
		processMap[proc.PID] = &processCopy
	}

	// Sort processes by CPU usage (descending)
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPUUsage > processes[j].CPUUsage
	})

	return processes, nil
}

// NewProcessInfo creates a new ProcessInfo for the given PID
func NewProcessInfo(pid int) (*ProcessInfo, error) {
	proc := &ProcessInfo{}
	proc.PID = pid
	err := proc.Read()
	return proc, err
}
