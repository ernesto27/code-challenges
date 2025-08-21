package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

type LoadAverage struct {
	load1min  float64
	load5min  float64
	load15min float64
}

var prevCPUStats *CPUStats
var processDisplayOffset int
var allProcesses []ProcessInfo
var processPageSize = 20

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
	prevCPUStats, _ = NewCPUStats()
	time.Sleep(100 * time.Millisecond) // Very brief pause
	printSystemInfo()

	// Create a ticker for 5-second intervals
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Channel for keyboard input
	inputCh := make(chan []byte, 1)

	// Goroutine to handle keyboard input
	go func() {
		for {
			var buf [3]byte
			n, err := os.Stdin.Read(buf[:])
			if err != nil {
				return
			}
			inputCh <- buf[:n]
		}
	}()

	// Main event loop
	for {
		select {
		case <-ticker.C:
			clearScreenAndHideCursor()
			printSystemInfo()

		case input := <-inputCh:
			// Handle keyboard input
			if len(input) == 1 && (input[0] == 'q' || input[0] == 'Q' || input[0] == '\x03') { // 'q', 'Q', or Ctrl+C
				fmt.Println()
				return
			}

			inputStr := string(input)
			switch inputStr {
			case "\x1b[B":
				if processDisplayOffset+processPageSize < len(allProcesses) {
					processDisplayOffset += processPageSize
					clearScreenAndHideCursor()
					printSystemInfo()
				}

			case "\x1b[A":
				if processDisplayOffset-processPageSize >= 0 {
					processDisplayOffset -= processPageSize
				} else {
					processDisplayOffset = 0
				}
				clearScreenAndHideCursor()
				printSystemInfo()
			}
		}

	}
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

func printSystemInfo() {
	currentTime := time.Now()
	cpuStats, err := NewCPUStats()

	var cpuBreakdown *CPUUsageBreakdown
	if err == nil {
		cpuBreakdown = cpuStats.CalculateUsage(prevCPUStats)
		prevCPUStats = cpuStats
	}

	memInfo, memErr := NewMemoryInfo()

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
		// Calculate used memory and components using the new methods
		usedKB := memInfo.GetUsedMemory()
		wiredKB := memInfo.GetWiredMemory()
		compressorKB := memInfo.GetCompressorMemory()
		unusedKB := memInfo.GetUnusedMemory()

		output.WriteString(fmt.Sprintf(" | PhysMem: %s used (%s wired, %s compressor), %s unused\n",
			memInfo.FormatSize(usedKB),
			memInfo.FormatSize(wiredKB),
			memInfo.FormatSize(compressorKB),
			memInfo.FormatSize(unusedKB)))
	} else {
		output.WriteString("\n")
	}

	output.WriteString("\n")
	output.WriteString("\r")

	// Get process data and add to output
	var err2 error
	allProcesses, err2 = getRunningProcess()
	if err2 == nil {
		output.WriteString("PID    COMMAND\n")

		// Calculate the range of processes to display
		startIdx := processDisplayOffset
		endIdx := processDisplayOffset + processPageSize
		if endIdx > len(allProcesses) {
			endIdx = len(allProcesses)
		}

		// Display only the current page of processes
		for i := startIdx; i < endIdx; i++ {
			proc := allProcesses[i]
			output.WriteString("\r")
			// Truncate command name if too long to maintain table formatting
			name := proc.Name
			if len(name) > 12 {
				name = name[:12]
			}
			output.WriteString(fmt.Sprintf("%-6d %s\n", proc.PID, name))
		}

		// Show pagination info
		totalProcesses := len(allProcesses)
		displayedEnd := endIdx
		if displayedEnd > totalProcesses {
			displayedEnd = totalProcesses
		}
		output.WriteString(fmt.Sprintf("\nShowing %d-%d of %d processes (Press ↓ for more)",
			startIdx+1, displayedEnd, totalProcesses))
	}

	// Clear screen, position cursor and print complete output
	fmt.Print("\033[2J\033[H")
	fmt.Print(output.String())
	os.Stdout.Sync()
}

type ProcessInfo struct {
	PID      int
	Name     string
	CPUUsage float64
}

func getRunningProcess() ([]ProcessInfo, error) {
	dataProcess, err := os.ReadDir("/proc/")
	if err != nil {
		return nil, err
	}

	var processes []ProcessInfo

	for _, entry := range dataProcess {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		commPath := filepath.Join("/proc", entry.Name(), "comm")
		data, err := os.ReadFile(commPath)
		if err != nil {
			continue
		}

		name := strings.TrimSpace(string(data))
		processes = append(processes, ProcessInfo{PID: pid, Name: name})
	}

	return processes, nil
}
