package main

import (
	"fmt"
	"os"
)

type RemoveComments struct {
	nameFile string
}

func newRemoveComments(nameFile string) *RemoveComments {
	code, err := os.ReadFile(nameFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	fmt.Println("File content:", string(code))
	return &RemoveComments{nameFile: nameFile}
}

func main() {

	rc := newRemoveComments("examples/python/calculator.py")
	if rc == nil {
		os.Exit(1)
	}
}
