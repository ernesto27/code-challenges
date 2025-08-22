package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type GitContribution struct {
	path string
}

func (g *GitContribution) getRepositories() []string {
	expandedPath := g.path
	if strings.HasPrefix(expandedPath, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			expandedPath = filepath.Join(home, expandedPath[2:])
		}
	}

	var repositories []string

	filepath.Walk(expandedPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() && info.Name() == ".git" {
			parentDir := filepath.Dir(path)
			repositories = append(repositories, parentDir)
			return filepath.SkipDir
		}

		return nil
	})

	return repositories
}

// getCurrentGitUser gets the current git user's email
func getCurrentGitUser() (string, error) {
	cmd := exec.Command("git", "config", "user.email")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting git user email: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// getCommitCount counts commits by the specified user in the last N days
func getCommitCount(repoPath string, userEmail string, days int) (int, error) {
	daysAgo := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	currentDir, err := os.Getwd()
	if err != nil {
		return 0, fmt.Errorf("error getting current directory: %w", err)
	}

	if err := os.Chdir(repoPath); err != nil {
		return 0, fmt.Errorf("error changing to repository directory: %w", err)
	}

	// Make sure to change back to the original directory when done
	defer os.Chdir(currentDir)

	cmd := exec.Command("git", "log", "--author="+userEmail, "--since="+daysAgo, "--oneline")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("error running git log: %w", err)
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		return 0, nil
	}

	commitLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return len(commitLines), nil
}

// getCommitDataByDate collects commit data from all repos for the given user and date range
func getCommitDataByDate(repositories []string, userEmail string, startDate time.Time) (map[string]int, error) {
	commitsByDate := make(map[string]int)

	// Get current directory to return to later
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}

	// Format the start date for git command
	startDateStr := startDate.Format("2006-01-02")

	// For each repository, get commit dates and count them
	for _, repoPath := range repositories {
		// Change to repository directory
		if err := os.Chdir(repoPath); err != nil {
			fmt.Printf("Warning: Could not access repository %s: %v\n", repoPath, err)
			continue
		}

		// Execute git log command to get commit timestamps
		cmd := exec.Command("git", "log", "--author="+userEmail, "--since="+startDateStr,
			"--format=%ad", "--date=short")
		output, err := cmd.Output()

		// Change back to original directory
		os.Chdir(currentDir)

		if err != nil {
			fmt.Printf("Warning: Error running git log in %s: %v\n", repoPath, err)
			continue
		}

		// Process the commit dates
		commitDates := strings.Split(strings.TrimSpace(string(output)), "\n")

		// Count commits by date
		for _, dateStr := range commitDates {
			if dateStr == "" {
				continue
			}
			commitsByDate[dateStr]++
		}
	}

	return commitsByDate, nil
}

func main() {
	// Define command line flags
	pathFlag := flag.String("path", "~/code", "Path to scan for git repositories")
	daysFlag := flag.Int("days", 30, "Number of days to analyze commit history")
	showGraphFlag := flag.Bool("graph", false, "Show contribution graph (uses hardcoded data for now)")
	flag.Parse()

	// If graph flag is set, show the graph and exit
	if *showGraphFlag {
		showContributionGraph()
		return
	}

	git := GitContribution{
		path: *pathFlag,
	}
	repositories := git.getRepositories()
	fmt.Println(repositories)

	userEmail, err := getCurrentGitUser()
	if err != nil {
		fmt.Printf("Error getting current git user: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Analyzing commits for user: %s in the past %d days\n", userEmail, *daysFlag)
	fmt.Println("----------------------------------------------------")

	if len(repositories) == 0 {
		fmt.Printf("No git repositories found in path: %s\n", *pathFlag)
		os.Exit(1)
	}

	totalCommits := 0
	for _, repo := range repositories {
		commitCount, err := getCommitCount(repo, userEmail, *daysFlag)
		if err != nil {
			fmt.Printf("Error counting commits in %s: %v\n", repo, err)
			continue
		}

		repoName := filepath.Base(repo)
		fmt.Printf("%-30s: %d commits\n", repoName, commitCount)
		totalCommits += commitCount
	}

	fmt.Println("----------------------------------------------------")
	fmt.Printf("Total commits in the last %d days: %d\n", *daysFlag, totalCommits)
}

// showContributionGraph displays a GitHub-style contribution graph with real git data
func showContributionGraph() {
	// Use the same path as regular analysis
	pathFlag := flag.Lookup("path").Value.String()

	git := GitContribution{
		path: pathFlag,
	}
	repositories := git.getRepositories()

	userEmail, err := getCurrentGitUser()
	if err != nil {
		fmt.Printf("Error getting current git user: %v\n", err)
		os.Exit(1)
	}

	// Create a contribution graph for the past year
	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0)
	// startDate := endDate.AddDate(0, 0, -30)

	fmt.Printf("Generating contribution graph for %s from %s to %s\n\n",
		userEmail, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// Create and initialize the graph
	graph := NewContributionGraph(startDate, endDate)

	if len(repositories) == 0 {
		fmt.Printf("No git repositories found in path: %s\n", pathFlag)
		os.Exit(1)
	}

	// Get commit data from repositories
	commitData, err := getCommitDataByDate(repositories, userEmail, startDate)
	if err != nil {
		fmt.Printf("Error getting commit data: %v\n", err)
		os.Exit(1)
	}

	// Use real data only
	graph.PopulateWithRealData(commitData)

	// Print the graph
	graph.PrintGraph()

	// Count total commits
	totalCommits := 0
	for _, day := range graph.Days {
		totalCommits += day.Contributions
	}

	fmt.Printf("\nTotal commits in the last year: %d\n", totalCommits)
}
