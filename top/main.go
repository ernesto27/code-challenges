package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table          table.Model
	processes      []*ProcessInfo
	processMap     map[int]*ProcessInfo
	cpuStats       *CPUStats
	coreUsage      map[int]float64
	lastUpdate     time.Time
	progressBars   map[int]progress.Model
	memoryInfo     *MemoryInfo
	memProgressBar progress.Model
	systemInfo     *SystemInfo
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
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
		data, err := m.cpuStats.CalculateUsagePerCore()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		m.coreUsage = data

		// Update memory info for progress bar
		_, err = m.memoryInfo.GetUsageInGB()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = m.systemInfo.GetUptime()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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
	// Build CPU cores header
	headerStr := m.buildCPUCoresHeader()

	// Build memory usage display
	memoryStr := m.buildMemoryHeader()

	// Build uptime display
	uptimeStr := m.buildUptimeHeader()

	tableStr := baseStyle.Render(m.table.View())

	return headerStr + "\n" + memoryStr + "\n" + uptimeStr + "\n\n" + tableStr + "\n\nPress q to quit • Use ↑/↓ to navigate"
}

func (m model) buildCPUCoresHeader() string {
	var coreStrs []string

	// Sort cores by index
	var coreIndices []int
	for coreIndex := range m.coreUsage {
		coreIndices = append(coreIndices, coreIndex)
	}
	sort.Ints(coreIndices)

	for _, coreIndex := range coreIndices {
		usage := m.coreUsage[coreIndex]

		// Get the progress bar for this core
		progBar, exists := m.progressBars[coreIndex]
		if !exists {
			continue
		}
		progressView := progBar.ViewAs(usage / 100.0)

		coreDisplay := fmt.Sprintf("CPU%d %s", coreIndex, progressView)

		// Add margin around each core display
		styledCore := lipgloss.NewStyle().
			MarginRight(2).
			MarginBottom(0).
			Render(coreDisplay)

		coreStrs = append(coreStrs, styledCore)
	}

	// Split cores into left and right columns
	var leftCol, rightCol []string

	midPoint := (len(coreStrs) + 1) / 2

	for i, coreStr := range coreStrs {
		if i < midPoint {
			leftCol = append(leftCol, coreStr)
		} else {
			rightCol = append(rightCol, coreStr)
		}
	}

	leftColumn := lipgloss.JoinVertical(lipgloss.Left, leftCol...)
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, rightCol...)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
}

func (m model) buildMemoryHeader() string {
	// Calculate percentage for progress bar
	percentage := m.memoryInfo.GetPercentageUse()
	progressView := m.memProgressBar.ViewAs(percentage)

	return fmt.Sprintf("Mem:%s %s", progressView, m.memoryInfo.usageInGB)
}

func (m model) buildUptimeHeader() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8BE9FD")).
		Bold(false)

	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#50FA7B")).
		Bold(true)

	label := labelStyle.Render("Uptime:")
	time := timeStyle.Render(m.systemInfo.Uptime)

	return fmt.Sprintf("%s %s", label, time)
}

func main() {
	// Initialize process data
	processMap := make(map[int]*ProcessInfo)

	// Take initial CPU reading
	cpuStats, err := NewCPUStats()
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
	t.SetStyles(s)
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

	// Initialize progress bars for CPU cores
	progressBars := make(map[int]progress.Model)

	// Calculate available width for progress bars
	// Get terminal width (default to 80 if unavailable)
	termWidth := 80
	if w, _, err := term.GetSize(0); err == nil {
		termWidth = w
	}

	// Calculate width per progress bar
	// Account for: core number (2), percentage (6), margins (6), 2 cores per row
	availableWidth := (termWidth - 20) / 2

	// Create progress bars for detected cores
	for i := 0; i < len(cpuStats.cores); i++ {
		prog := progress.New(
			progress.WithScaledGradient("#00ff00", "#ff0000"),
			progress.WithWidth(availableWidth),
		)
		progressBars[i] = prog
	}

	// Create memory progress bar
	memProgressBar := progress.New(
		progress.WithScaledGradient("#00ff00", "#ff0000"),
		progress.WithWidth(availableWidth),
	)

	// Initialize memory info
	memoryInfo := NewMemoryInfo()

	systemInfo := NewSystemInfo()

	// Create model
	m := model{
		table:          t,
		processMap:     processMap,
		cpuStats:       cpuStats,
		lastUpdate:     time.Now(),
		progressBars:   progressBars,
		memoryInfo:     memoryInfo,
		memProgressBar: memProgressBar,
		systemInfo:     systemInfo,
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
