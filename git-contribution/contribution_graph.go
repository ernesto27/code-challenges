package main

import (
	"fmt"
	"time"
)

// ContributionLevel represents intensity of contributions
type ContributionLevel int

const (
	NoContribution ContributionLevel = iota
	LowContribution
	MediumContribution
	HighContribution
	VeryHighContribution
)

// ContributionDay represents a single day with contribution data
type ContributionDay struct {
	Date          time.Time
	Contributions int
	Level         ContributionLevel
}

// ContributionGraph represents a full year of contribution data
type ContributionGraph struct {
	Days      []ContributionDay
	StartDate time.Time
	EndDate   time.Time
}

// NewContributionGraph creates a new contribution graph with the given time range
func NewContributionGraph(startDate, endDate time.Time) *ContributionGraph {
	days := []ContributionDay{}

	// Generate days between start and end dates
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		days = append(days, ContributionDay{
			Date:          d,
			Contributions: 0,
			Level:         NoContribution,
		})
	}

	return &ContributionGraph{
		Days:      days,
		StartDate: startDate,
		EndDate:   endDate,
	}
}

// PopulateWithHardcodedData fills the contribution graph with hardcoded sample data
func (g *ContributionGraph) PopulateWithHardcodedData() {
	// Create a deterministic pattern that looks similar to the image
	for i := range g.Days {
		day := &g.Days[i]

		// Assign contributions based on patterns
		switch {
		case day.Date.Weekday() == time.Monday && i%7 == 0:
			day.Contributions = 4
			day.Level = MediumContribution
		case day.Date.Weekday() == time.Friday:
			if i%3 == 0 {
				day.Contributions = 8
				day.Level = HighContribution
			}
		case day.Date.Month() == time.February && day.Date.Day()%5 == 0:
			day.Contributions = 12
			day.Level = VeryHighContribution
		case day.Date.Month() == time.October && day.Date.Day()%4 == 0:
			day.Contributions = 6
			day.Level = MediumContribution
		case day.Date.Weekday() == time.Wednesday && i%2 == 0:
			day.Contributions = 2
			day.Level = LowContribution
		}

		// Ensure some high activity streaks
		if i > 5 && i < 20 && i%3 == 0 {
			day.Contributions = 10
			day.Level = HighContribution
		}

		// Add some random high days
		if (i+3)%17 == 0 {
			day.Contributions = 15
			day.Level = VeryHighContribution
		}
	}
}

// PopulateWithRealData fills the contribution graph with actual git commit data
func (g *ContributionGraph) PopulateWithRealData(commitsByDate map[string]int) {
	// Find max commits in a day to calculate relative levels
	maxCommits := 1 // Default to at least 1 to avoid division by zero
	for _, count := range commitsByDate {
		if count > maxCommits {
			maxCommits = count
		}
	}

	// Set contribution levels based on commit counts
	for i := range g.Days {
		dateStr := g.Days[i].Date.Format("2006-01-02")
		commits := commitsByDate[dateStr]

		// Set the commits count
		g.Days[i].Contributions = commits

		// Set the level based on number of commits relative to max
		if commits == 0 {
			g.Days[i].Level = NoContribution
		} else if float64(commits) <= float64(maxCommits)*0.25 {
			g.Days[i].Level = LowContribution
		} else if float64(commits) <= float64(maxCommits)*0.5 {
			g.Days[i].Level = MediumContribution
		} else if float64(commits) <= float64(maxCommits)*0.75 {
			g.Days[i].Level = HighContribution
		} else {
			g.Days[i].Level = VeryHighContribution
		}
	}
}

// getLevelColor returns ANSI color code for the contribution level
func getLevelColor(level ContributionLevel) string {
	switch level {
	case NoContribution:
		return "\033[48;5;235m" // Dark gray background
	case LowContribution:
		return "\033[48;5;22m" // Light green background
	case MediumContribution:
		return "\033[48;5;28m" // Medium green background
	case HighContribution:
		return "\033[48;5;34m" // Strong green background
	case VeryHighContribution:
		return "\033[48;5;40m" // Very bright green background
	default:
		return "\033[0m" // Reset
	}
}

// PrintGraph prints the contribution graph to the console
func (g *ContributionGraph) PrintGraph() {
	resetColor := "\033[0m"
	square := "  " // Two spaces for a square

	// Print months
	fmt.Print("     ") // Offset for weekday labels
	currentMonth := g.Days[0].Date.Month()
	for _, day := range g.Days {
		if day.Date.Month() != currentMonth {
			currentMonth = day.Date.Month()
			fmt.Printf("%-8s", day.Date.Month().String()[:3])
		}
	}
	fmt.Println()

	// Group by weekday
	weekdaysDisplay := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	for weekday := 0; weekday < 7; weekday++ {
		// Print weekday label every other day
		if weekday%2 == 1 || weekday == 0 || weekday == 6 {
			fmt.Printf("%-4s ", weekdaysDisplay[weekday])
		} else {
			fmt.Print("     ")
		}

		// Print all days for this weekday
		for _, day := range g.Days {
			if int(day.Date.Weekday()) == weekday {
				color := getLevelColor(day.Level)
				fmt.Print(color + square + resetColor)
			}
		}
		fmt.Println()
	}

	// Print legend
	fmt.Println("\nLess",
		getLevelColor(NoContribution)+square+resetColor,
		getLevelColor(LowContribution)+square+resetColor,
		getLevelColor(MediumContribution)+square+resetColor,
		getLevelColor(HighContribution)+square+resetColor,
		getLevelColor(VeryHighContribution)+square+resetColor,
		"More")
}
