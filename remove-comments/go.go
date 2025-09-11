package main

import (
	"fmt"
	"os"
	"strings"
)

type RemoveCommentsGo struct {
	RemoveCommentsC
}

func newRemoveCommentsGo(nameFile string) *RemoveCommentsGo {
	code, err := os.ReadFile(nameFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	rc := &RemoveCommentsC{
		nameFile:        nameFile,
		originalContent: string(code),
	}

	return &RemoveCommentsGo{
		RemoveCommentsC: *rc,
	}
}

func (rg *RemoveCommentsGo) Remove() string {
	lines := strings.Split(rg.originalContent, "\n")
	var result []string
	inBlockComment := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		trimLine := strings.TrimSpace(line)
		fmt.Println(trimLine)
		fmt.Println(strings.HasPrefix(trimLine, "//go:"))

		if strings.HasPrefix(trimLine, "//go:") {
			result = append(result, trimLine)
			continue
		}

		processedLine, stillInBlock := rg.processLine(line, inBlockComment)
		inBlockComment = stillInBlock

		if trimLine != "" {
			result = append(result, processedLine)
		} else {
			if trimLine == "" && !inBlockComment {
				result = append(result, "")
			}
		}
	}

	return strings.Join(result, "\n")
}
