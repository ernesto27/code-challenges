package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const redColor = "\033[31m"
const resetColor = "\033[0m"

type Grep struct {
}

func (g *Grep) Run(filter string, content string, exclude bool, caseInsensitive bool) (string, error) {
	if filter == "" {
		return content, nil
	}

	contentLines := strings.Split(content, "\n")
	var result []string

	if caseInsensitive {
		filter = "(?i)" + filter
	}

	re, err := regexp.Compile(filter)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return "", err
	}

	for _, line := range contentLines {
		if exclude {
			if !re.MatchString(line) {
				result = append(result, line)
			}

		} else {
			coloredLine := re.ReplaceAllStringFunc(line, func(match string) string {
				return redColor + match + resetColor
			})

			if coloredLine != line {
				result = append(result, coloredLine)
			}
		}

	}

	return strings.Join(result, "\n"), nil
}

func (g *Grep) RunRecursive(filter string, initialPath string) (string, error) {
	var results []string

	err := filepath.Walk(initialPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			filteredContent, err := g.Run(filter, string(content), false, false)
			if err != nil {
				return err
			}

			if filteredContent != "" {
				filteredContentLines := strings.Split(filteredContent, "\n")

				for _, line := range filteredContentLines {
					results = append(results, fmt.Sprintf("%s:%s", path, line))
				}
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return strings.Join(results, "\n"), nil
}

func main() {
	recursive := flag.Bool("r", false, "recursive search")
	exclude := flag.Bool("v", false, "exclude response")
	caseInsensitive := flag.Bool("i", false, "case insensitive search")
	flag.Parse()

	if len(os.Args) < 3 {
		fmt.Println("Usage: grep <filter> <file>")
		return
	}

	grep := Grep{}

	var resp string
	var err error

	if *recursive {
		initialPath := os.Args[3]
		if os.Args[3] == "*" {
			initialPath = "."
		}

		fmt.Println("Initial path:", initialPath)

		resp, err = grep.RunRecursive(os.Args[2], initialPath)
		if err != nil {
			fmt.Println("Error running grep:", err)
			return
		}
	} else {
		filter := os.Args[1]
		filePath := os.Args[2]
		if *exclude {
			filter = os.Args[2]
			filePath = os.Args[3]
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		content := string(data)
		resp, err = grep.Run(filter, content, *exclude, *caseInsensitive)
		if err != nil {
			fmt.Println("Error running grep:", err)
			return
		}
	}

	fmt.Println(resp)

}
