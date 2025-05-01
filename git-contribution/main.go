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

func main() {
	// Define command line flags
	pathFlag := flag.String("path", "~/code", "Path to scan for git repositories")
	daysFlag := flag.Int("days", 30, "Number of days to analyze commit history")
	flag.Parse()

	git := GitContribution{
		path: *pathFlag,
	}
	repositories := git.getRepositories()

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
