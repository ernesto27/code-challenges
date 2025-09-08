package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type RemoveCommentsInterface interface {
	Remove() string
}

func main() {
	//filePath := "examples/python/calculator.py"
	//filePath := "examples/javascript/calculator.js"
	filePath := "examples/c/calculator.c"

	extension := filepath.Ext(filePath)

	var rc RemoveCommentsInterface
	switch extension {
	case ".py":
		rc = newRemoveCommentsPython(filePath)
	case ".js":
		rc = newRemoveCommentsC(filePath)
	case ".c":
		rc = newRemoveCommentsC(filePath)
	default:
		panic("Unsupported file type")
	}

	// Process the content using method
	resp := rc.Remove()

	fmt.Println("File without comments and docstrings:\n", resp)

	if rc == nil {
		os.Exit(1)
	}
}
