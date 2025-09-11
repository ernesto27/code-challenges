package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RemoveCommentsInterface interface {
	Remove() string
}

func saveCleanedCode(cleanedCode, originalFilePath string) error {
	dir := filepath.Dir(originalFilePath)
	filename := filepath.Base(originalFilePath)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	cleanedFilePath := filepath.Join(dir, nameWithoutExt+"_cleaned"+ext)

	err := os.WriteFile(cleanedFilePath, []byte(cleanedCode), 0644)
	if err != nil {
		return fmt.Errorf("error saving cleaned code: %v", err)
	}

	fmt.Printf("Cleaned code saved to: %s\n", cleanedFilePath)
	return nil
}

func main() {
	//filePath := "examples/python/calculator.py"
	//filePath := "examples/javascript/calculator.js"
	//filePath := "examples/c/calculator.c"
	filePath := "examples/go/calculator.go"

	extension := filepath.Ext(filePath)

	var rc RemoveCommentsInterface
	switch extension {
	case ".py":
		rc = newRemoveCommentsPython(filePath)
	case ".js":
		rc = newRemoveCommentsC(filePath)
	case ".c":
		rc = newRemoveCommentsC(filePath)
	case ".go":
		rc = newRemoveCommentsGo(filePath)
	default:
		panic("Unsupported file type")
	}

	// Process the content using method
	resp := rc.Remove()

	fmt.Println("File without comments and docstrings:\n", resp)

	// Save cleaned code to disk
	// err := saveCleanedCode(resp, filePath)
	// if err != nil {
	// 	fmt.Printf("Failed to save cleaned code: %v\n", err)
	// 	os.Exit(1)
	// }

	if rc == nil {
		os.Exit(1)
	}
}
