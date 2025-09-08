package main

import (
	"fmt"
	"os"
	"strings"
)

type RemoveCommentsC struct {
	nameFile        string
	originalContent string
	cleanedContent  string
}

func newRemoveCommentsC(nameFile string) *RemoveCommentsC {
	code, err := os.ReadFile(nameFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	return &RemoveCommentsC{
		nameFile:        nameFile,
		originalContent: string(code),
	}
}

func (rc *RemoveCommentsC) Remove() string {
	lines := strings.Split(rc.originalContent, "\n")
	var result []string
	inBlockComment := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		processedLine, stillInBlock := rc.processLine(line, inBlockComment)
		inBlockComment = stillInBlock

		if strings.TrimSpace(processedLine) != "" {
			result = append(result, processedLine)
		} else {
			if strings.TrimSpace(line) == "" && !inBlockComment {
				result = append(result, "")
			}
		}
	}

	return strings.Join(result, "\n")
}

func (rc *RemoveCommentsC) processLine(line string, inBlockComment bool) (string, bool) {
	if inBlockComment {
		if endPos := strings.Index(line, "*/"); endPos != -1 {
			remaining := line[endPos+2:]
			// Trim leading whitespace after block comment end
			remaining = strings.TrimLeft(remaining, " \t")
			return rc.processLine(remaining, false)
		}
		return "", true
	}

	var result strings.Builder
	i := 0

	for i < len(line) {
		char := line[i]
		// fmt.Println("Char:", string(char))

		// Check for single-line comment
		if i < len(line)-1 && line[i:i+2] == "//" {
			if !rc.isInsideString(line, i) {
				// Remove trailing whitespace before comment
				resultStr := strings.TrimRight(result.String(), " \t")
				return resultStr, false
			}
		}

		// Check for block comment
		if i < len(line)-1 && line[i:i+2] == "/*" {
			if !rc.isInsideString(line, i) {
				if endPos := strings.Index(line[i+2:], "*/"); endPos != -1 {
					// Block comment ends on same line, skip it entirelycleanedContent
					i = i + 2 + endPos + 2
					// Skip any whitespace immediately after the block comment
					for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
						i++
					}
					continue
				} else {
					// Block comment continues to next line
					resultStr := strings.TrimRight(result.String(), " \t")
					return resultStr, true
				}
			}
		}

		result.WriteByte(char)
		i++
	}

	return result.String(), false
}

func (rc *RemoveCommentsC) isInsideString(line string, position int) bool {
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false
	inTemplateString := false

	for i := 0; i < position; i++ {
		char := line[i]
		// fmt.Println("isInsideString Char:", string(char))

		if escaped {
			escaped = false
			continue
		}

		if char == '\\' && (inSingleQuote || inDoubleQuote || inTemplateString) {
			escaped = true
			continue
		}

		switch char {
		case '\'':
			if !inDoubleQuote && !inTemplateString {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote && !inTemplateString {
				inDoubleQuote = !inDoubleQuote
			}
		case '`':
			if !inSingleQuote && !inDoubleQuote {
				inTemplateString = !inTemplateString
			}
		}
	}

	return inSingleQuote || inDoubleQuote || inTemplateString
}
