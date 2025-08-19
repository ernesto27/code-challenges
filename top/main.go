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
	user, nice, system, idle, iowait, irq, softirq, steal uint64
}

var prevCPUStats *CPUStats

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
	defer term.Restore(fd, oldState)

	// Put the terminal into raw mode. This disables input echoing.
	_, err = term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}

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

func calculateCPUUsage(prev, curr *CPUStats) float64 {
	if prev == nil {
		return 0.0
	}

	prevTotal := prev.user + prev.nice + prev.system + prev.idle + prev.iowait + prev.irq + prev.softirq + prev.steal
	currTotal := curr.user + curr.nice + curr.system + curr.idle + curr.iowait + curr.irq + curr.softirq + curr.steal

	totalDiff := currTotal - prevTotal
	idleDiff := curr.idle - prev.idle

	if totalDiff == 0 {
		return 0.0
	}

	return (1.0 - float64(idleDiff)/float64(totalDiff)) * 100.0
}

func printSystemInfo() {
	currentTime := time.Now()
	cpuStats, err := readCPUStats()

	var cpuUsage float64
	if err == nil {
		cpuUsage = calculateCPUUsage(prevCPUStats, cpuStats)
		prevCPUStats = cpuStats
	}

	if err != nil {
		fmt.Printf("\rTime: %s | CPU: --%%", currentTime.Format("15:04:05"))
	} else {
		fmt.Printf("\rTime: %s | CPU: %.1f%%", currentTime.Format("15:04:05"), cpuUsage)
	}
}
