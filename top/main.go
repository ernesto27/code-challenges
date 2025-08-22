package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table        table.Model
	processes    []*ProcessInfo
	processMap   map[int]*ProcessInfo
	cpuStats     *CPUStats
	prevCPUStats *CPUStats
	lastUpdate   time.Time
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tickMsg:
		m.updateProcessData()
		return m, tickCmd()
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *model) updateProcessData() {
	// Get new process data
	newProcesses, err := GetRunningProcesses()
	if err != nil {
		return
	}

	// Update CPU calculations
	for _, proc := range newProcesses {
		if prevProc, exists := m.processMap[proc.PID]; exists {
			// Copy previous measurements
			proc.PrevUTime = prevProc.UTime
			proc.PrevSTime = prevProc.STime
			proc.PrevMeasurementTime = prevProc.PrevMeasurementTime
			proc.HasPrevMeasurement = prevProc.HasPrevMeasurement
			proc.CalculateCPUUsage()
		}
		proc.UpdateMeasurement()
		processCopy := *proc
		m.processMap[proc.PID] = &processCopy
	}

	// Sort processes by CPU usage (descending)
	sort.Slice(newProcesses, func(i, j int) bool {
		return newProcesses[i].CPUUsage > newProcesses[j].CPUUsage
	})

	m.processes = newProcesses

	// Convert processes to table rows
	rows := make([]table.Row, len(m.processes))
	for i, proc := range m.processes {
		name := proc.Name
		if len(name) > 15 {
			name = name[:15]
		}

		rows[i] = table.Row{
			strconv.Itoa(proc.PID),
			fmt.Sprintf("%.1f", proc.CPUUsage),
			name,
			proc.GetTime(),
			strconv.Itoa(proc.ThreadsCount),
			proc.Username,
		}
	}

	m.table.SetRows(rows)
	m.lastUpdate = time.Now()
}

func (m model) View() string {
	// System info header
	currentTime := time.Now()

	var header strings.Builder
	header.WriteString(fmt.Sprintf("Time: %s", currentTime.Format("15:04:05")))

	// Get CPU info
	if cpuStats, err := NewCPUStats(); err == nil {
		if m.prevCPUStats != nil {
			cpuBreakdown := cpuStats.CalculateUsage(m.prevCPUStats)
			header.WriteString(fmt.Sprintf(" | CPU: %.1f%% user, %.1f%% sys, %.1f%% idle",
				cpuBreakdown.userPercent, cpuBreakdown.systemPercent, cpuBreakdown.idlePercent))
		}
	}

	// Get memory info
	if memInfo, err := NewMemoryInfo(); err == nil && memInfo.totalKB > 0 {
		usedKB := memInfo.GetUsedMemory()
		unusedKB := memInfo.GetUnusedMemory()
		header.WriteString(fmt.Sprintf(" | PhysMem: %s used, %s unused",
			memInfo.FormatSize(usedKB), memInfo.FormatSize(unusedKB)))
	}

	headerStr := header.String()
	tableStr := baseStyle.Render(m.table.View())

	return headerStr + "\n\n" + tableStr + "\n\nPress q to quit • Use ↑/↓ to navigate"
}

func main() {
	// Initialize process data
	processMap := make(map[int]*ProcessInfo)

	// Take initial CPU reading
	prevCPUStats, err := NewCPUStats()
	if err != nil {
		fmt.Printf("Error reading CPU stats: %v\n", err)
		os.Exit(1)
	}
	time.Sleep(100 * time.Millisecond)

	// Define table columns
	columns := []table.Column{
		{Title: "PID", Width: 8},
		{Title: "%CPU", Width: 6},
		{Title: "COMMAND", Width: 15},
		{Title: "TIME+", Width: 9},
		{Title: "#TH", Width: 4},
		{Title: "USER", Width: 12},
	}

	// Create initial empty table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	// Style the table
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	// Create model
	m := model{
		table:        t,
		processMap:   processMap,
		prevCPUStats: prevCPUStats,
		lastUpdate:   time.Now(),
	}

	// Initial data load
	m.updateProcessData()

	// Run the program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
