package main

import (
	"fmt"
	"os"
	"strings"
)

type RemoveComments struct {
	nameFile        string
	originalContent string
	cleanedContent  string
}

func newRemoveComments(nameFile string) *RemoveComments {
	code, err := os.ReadFile(nameFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	rc := &RemoveComments{
		nameFile:        nameFile,
		originalContent: string(code),
	}
	
	// Process the content using method
	rc.cleanedContent = rc.removeCommentsAndDocstrings(rc.originalContent)

	fmt.Println("File without comments and docstrings:\n", rc.cleanedContent)
	return rc
}

func (rc *RemoveComments) isInsideString(line string, position int) bool {
	inSingleQuote := false
	inDoubleQuote := false
	inTripleDouble := false
	inTripleSingle := false
	escaped := false

	for i := 0; i < position; i++ {
		char := line[i]

		if escaped {
			escaped = false
			continue
		}

		if char == '\\' && !inTripleDouble && !inTripleSingle {
			escaped = true
			continue
		}

		if i+3 <= len(line) && line[i:i+3] == `"""` && !inSingleQuote && !inTripleSingle {
			inTripleDouble = !inTripleDouble
			i += 2
			continue
		}

		if i+3 <= len(line) && line[i:i+3] == `'''` && !inDoubleQuote && !inTripleDouble {
			inTripleSingle = !inTripleSingle
			i += 2
			continue
		}

		if char == '\'' && !inDoubleQuote && !inTripleDouble && !inTripleSingle {
			inSingleQuote = !inSingleQuote
		} else if char == '"' && !inSingleQuote && !inTripleDouble && !inTripleSingle {
			inDoubleQuote = !inDoubleQuote
		}
	}

	return inSingleQuote || inDoubleQuote || inTripleDouble || inTripleSingle
}

func (rc *RemoveComments) removeCommentsFromLine(line string) (string, bool) {
	trimmed := strings.TrimSpace(line) 
	if strings.HasPrefix(trimmed, "#") {
		return "", false 
	}

	for i := 0; i < len(line); i++ {
		if line[i] == '#' {
			if !rc.isInsideString(line, i) {
				clened := strings.TrimRight(line[:i], " \t")
				return clened, true 
			}
		}
	}

	return line, true
}

func (rc *RemoveComments) isDocstring(line string, prevLines []string) bool {
	trimmed := strings.TrimSpace(line)
	
	// Check if line starts with triple quotes
	if !strings.HasPrefix(trimmed, `"""`) && !strings.HasPrefix(trimmed, `'''`) {
		return false
	}
	
	// Look at previous non-empty lines to determine context
	for i := len(prevLines) - 1; i >= 0; i-- {
		prev := strings.TrimSpace(prevLines[i])
		if prev == "" {
			continue // Skip empty lines
		}
		
		// Check if previous line indicates start of function/class/module
		if strings.HasPrefix(prev, "def ") ||
		   strings.HasPrefix(prev, "class ") ||
		   strings.HasSuffix(prev, ":") ||
		   i == 0 { // First line of file (module docstring)
			return true
		}
		
		// If we hit a non-definition line, it's not a docstring
		return false
	}
	
	// If no previous lines, could be module docstring
	return len(prevLines) == 0
}

func (rc *RemoveComments) skipDocstring(lines []string, startIndex int) int {
	line := strings.TrimSpace(lines[startIndex])
	
	// Determine quote type
	var quoteType string
	if strings.HasPrefix(line, `"""`) {
		quoteType = `"""`
	} else if strings.HasPrefix(line, `'''`) {
		quoteType = `'''`
	} else {
		return startIndex // Not a docstring start
	}
	
	// Check if it's a single-line docstring
	if strings.Count(line, quoteType) >= 2 && len(line) > len(quoteType) {
		// Single-line docstring like """This is a docstring"""
		return startIndex
	}
	
	// Multi-line docstring - find the closing quotes
	for i := startIndex + 1; i < len(lines); i++ {
		if strings.Contains(lines[i], quoteType) {
			return i
		}
	}
	
	// If we reach here, docstring wasn't properly closed
	return startIndex
}

func (rc *RemoveComments) removeCommentsAndDocstrings(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	var processedLines []string
	
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		
		// First remove comments from the line
		cleanedLine, keepLine := rc.removeCommentsFromLine(line)

		if !keepLine {
			// Entire line was a comment, skip it
			continue
		}
		
		// Check if this line starts a docstring
		if rc.isDocstring(cleanedLine, processedLines) {
			// Skip the docstring
			endIndex := rc.skipDocstring(lines, i)
			i = endIndex
			continue
		}
					
		result = append(result, cleanedLine)
		processedLines = append(processedLines, cleanedLine)
	}
	
	return strings.Join(result, "\n")
}

// Step 6: Output Methods
func (rc *RemoveComments) GetCleanedContent() string {
	return rc.cleanedContent
}

func (rc *RemoveComments) GetOriginalContent() string {
	return rc.originalContent
}

func (rc *RemoveComments) SaveToFile(outputPath string) error {
	return os.WriteFile(outputPath, []byte(rc.cleanedContent), 0644)
}

func (rc *RemoveComments) ReplaceOriginal() error {
	return rc.SaveToFile(rc.nameFile)
}

func main() {
	rc := newRemoveComments("examples/python/calculator.py")
	if rc == nil {
		os.Exit(1)
	}
}
