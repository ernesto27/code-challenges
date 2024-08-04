package main

import (
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

func (g *Grep) Run(filter string, content string, exclude bool) (string, error) {
	if filter == "" {
		return content, nil
	}

	contentLines := strings.Split(content, "\n")
	var result []string

	// re, err := regexp.Compile("(?i)" + filter)
	//excludePattern := fmt.Sprintf("^(?!.*%s).*", regexp.QuoteMeta(filter))

	re, err := regexp.Compile(filter)
	//re, err := regexp.Compile(excludePattern)
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

			filteredContent, err := g.Run(filter, string(content), false)
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
	// if len(os.Args) < 3 {
	// 	fmt.Println("Usage: grep <filter> <file>")
	// 	return
	// }

	data, err := os.ReadFile("custom.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	content := string(data)

	grep := Grep{}
	resp, err := grep.Run("Nirvana", content, true)
	if err != nil {
		fmt.Println("Error running grep:", err)
		return
	}

	fmt.Println(resp)

	// return

	// data, err := os.ReadFile("rockbands.txt")
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }

	// grep := Grep{}
	// resp, err := grep.Run("j", string(data))
	// if err != nil {
	// 	fmt.Println("Error running grep:", err)
	// 	return
	// }

	//fmt.Println(resp)

	// grep := Grep{}
	// resp, err := grep.RunRecursive("Madonna", ".")
	// if err != nil {
	// 	fmt.Println("Error running grep:", err)
	// 	return
	// }
	// fmt.Println(resp)

}
