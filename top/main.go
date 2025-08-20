package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
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

type LoadAverage struct {
	load1min  float64
	load5min  float64
	load15min float64
}

var prevCPUStats *CPUStats

// clearScreenAndHideCursor clears the terminal screen and hides the cursor
func clearScreenAndHideCursor() {
	fmt.Print("\033[2J\033[H\033[?25l")
}

// showCursor shows the terminal cursor
func showCursor() {
	fmt.Print("\033[?25h")
}

func main() {
	// Get the file descriptor for standard input
	fd := int(os.Stdin.Fd())

	// Check if standard input is a terminal
	if !term.IsTerminal(fd) {
		fmt.Println("Not running in a terminal.")
		return
	}

	// Save the original terminal state so we can restore it later
	oldState, err := term.GetState(fd)
	if err != nil {
		panic(err)
	}
	defer func() {
		showCursor()
		term.Restore(fd, oldState)
	}()

	// Put the terminal into raw mode. This disables input echoing.
	_, err = term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}

	// Clear the terminal screen and hide cursor
	clearScreenAndHideCursor()

	// Take two quick readings to calculate initial CPU usage
	prevCPUStats, _ = readCPUStats()
	time.Sleep(100 * time.Millisecond) // Very brief pause
	printSystemInfo()

	// Create a ticker for 5-second intervals
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Channel for keyboard input
	inputCh := make(chan byte, 1)

	// Goroutine to handle keyboard input
	go func() {
		for {
			var buf [1]byte
			_, err := os.Stdin.Read(buf[:])
			if err != nil {
				return
			}
			inputCh <- buf[0]
		}
	}()

	// Main event loop
	for {
		select {
		case <-ticker.C:
			clearScreenAndHideCursor()
			printSystemInfo()

		case char := <-inputCh:
			// Handle keyboard input
			if char == 'q' || char == 'Q' || char == 3 { // 'q', 'Q', or Ctrl+C
				fmt.Println()
				return
			}
		}
	}
}

func readCPUStats() (*CPUStats, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read CPU line")
	}

	line := scanner.Text()
	fields := strings.Fields(line)
	if len(fields) < 8 || fields[0] != "cpu" {
		return nil, fmt.Errorf("invalid CPU line format")
	}

	stats := &CPUStats{}
	values := []*uint64{&stats.user, &stats.nice, &stats.system, &stats.idle,
		&stats.iowait, &stats.irq, &stats.softirq, &stats.steal}

	for i, val := range values {
		parsed, err := strconv.ParseUint(fields[i+1], 10, 64)
		if err != nil {
			return nil, err
		}
		*val = parsed
	}

	return stats, nil
}

type CPUUsageBreakdown struct {
	userPercent   float64
	systemPercent float64
	idlePercent   float64
}

func calculateCPUBreakdown(prev, curr *CPUStats) *CPUUsageBreakdown {
	if prev == nil {
		return &CPUUsageBreakdown{0.0, 0.0, 0.0}
	}

	prevTotal := prev.user + prev.nice + prev.system + prev.idle + prev.iowait + prev.irq + prev.softirq + prev.steal
	currTotal := curr.user + curr.nice + curr.system + curr.idle + curr.iowait + curr.irq + curr.softirq + curr.steal

	totalDiff := currTotal - prevTotal
	if totalDiff == 0 {
		return &CPUUsageBreakdown{0.0, 0.0, 0.0}
	}

	userDiff := (curr.user + curr.nice) - (prev.user + prev.nice)
	systemDiff := (curr.system + curr.irq + curr.softirq) - (prev.system + prev.irq + prev.softirq)
	idleDiff := curr.idle - prev.idle

	return &CPUUsageBreakdown{
		userPercent:   (float64(userDiff) / float64(totalDiff)) * 100.0,
		systemPercent: (float64(systemDiff) / float64(totalDiff)) * 100.0,
		idlePercent:   (float64(idleDiff) / float64(totalDiff)) * 100.0,
	}
}

func readMemoryInfo() (*MemoryInfo, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	memInfo := &MemoryInfo{}
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
			memInfo.totalKB = value
			fieldsFound++
		case "MemAvailable":
			memInfo.availableKB = value
			fieldsFound++
		case "MemFree":
			memInfo.freeKB = value
			fieldsFound++
		case "Buffers":
			memInfo.buffersKB = value
			fieldsFound++
		case "Cached":
			memInfo.cachedKB = value
			fieldsFound++
		case "Unevictable":
			memInfo.unevictableKB = value
			fieldsFound++
		case "Active":
			memInfo.activeKB = value
			fieldsFound++
		case "SwapCached":
			memInfo.swapCachedKB = value
			fieldsFound++
		}

		// Stop early if we have all target values
		if fieldsFound >= targetFields {
			break
		}
	}

	return memInfo, scanner.Err()
}

func readLoadAverage() (*LoadAverage, error) {
	file, err := os.Open("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read loadavg line")
	}

	line := scanner.Text()
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return nil, fmt.Errorf("invalid loadavg format")
	}

	loadAvg := &LoadAverage{}

	load1, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 1min load average: %v", err)
	}
	loadAvg.load1min = load1

	load5, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 5min load average: %v", err)
	}
	loadAvg.load5min = load5

	load15, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 15min load average: %v", err)
	}
	loadAvg.load15min = load15

	return loadAvg, nil
}

func formatMemorySize(kb uint64) string {
	if kb >= 1024*1024 {
		gb := float64(kb) / (1024 * 1024)
		return fmt.Sprintf("%.0fG", gb)
	} else if kb >= 1024 {
		mb := float64(kb) / 1024
		return fmt.Sprintf("%.0fM", mb)
	}
	return fmt.Sprintf("%dK", kb)
}

func printSystemInfo() {
	currentTime := time.Now()
	cpuStats, err := readCPUStats()

	var cpuBreakdown *CPUUsageBreakdown
	if err == nil {
		cpuBreakdown = calculateCPUBreakdown(prevCPUStats, cpuStats)
		prevCPUStats = cpuStats
	}

	memInfo, memErr := readMemoryInfo()

	loadAvg, loadErr := readLoadAverage()

	// Build the complete output string first
	var output strings.Builder
	output.WriteString(fmt.Sprintf("Time: %s | ", currentTime.Format("15:04:05")))

	if loadErr != nil {
		output.WriteString("Load Avg: --, --, -- | ")
	} else {
		output.WriteString(fmt.Sprintf("Load Avg: %.2f, %.2f, %.2f | ", loadAvg.load1min, loadAvg.load5min, loadAvg.load15min))
	}

	if err != nil {
		output.WriteString("CPU usage: --% user, --% sys, --% idle")
	} else {
		output.WriteString(fmt.Sprintf("CPU usage: %.2f%% user, %.1f%% sys, %.2f%% idle", cpuBreakdown.userPercent, cpuBreakdown.systemPercent, cpuBreakdown.idlePercent))
	}

	if memErr == nil && memInfo.totalKB > 0 {
		// Calculate used memory and components
		usedKB := memInfo.totalKB - memInfo.availableKB
		wiredKB := memInfo.unevictableKB + (memInfo.activeKB / 4)                     // Approximation for "wired" memory
		compressorKB := memInfo.buffersKB + memInfo.swapCachedKB + memInfo.cachedKB/2 // Approximation for "compressor"
		unusedKB := memInfo.freeKB

		output.WriteString(fmt.Sprintf(" | PhysMem: %s used (%s wired, %s compressor), %s unused",
			formatMemorySize(usedKB),
			formatMemorySize(wiredKB),
			formatMemorySize(compressorKB),
			formatMemorySize(unusedKB)))
	}

	// Clear line and print the complete output on same line
	fmt.Printf("\r%s", output.String())
	os.Stdout.Sync()
}
