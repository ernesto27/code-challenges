package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func main() {
	git := GitContribution{
		path: "~/code",
	}
	repositories := git.getRepositories()

	for _, repo := range repositories {
		fmt.Println(repo)
	}
}
