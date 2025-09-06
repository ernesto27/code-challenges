package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-python/gpython/parser"
	"github.com/go-python/gpython/py"
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

	pythonCodeStr := string(code)
	astRootNode, err := parser.Parse(strings.NewReader(pythonCodeStr), nameFile, py.ExecMode)
	if err != nil {
		log.Fatalf("Error parsing Python code: %v", err)
	}

	fmt.Printf("AST Root Node: %+v\n", astRootNode)

	return &RemoveComments{nameFile: nameFile}
}

func main() {

	rc := newRemoveComments("examples/python/calculator.py")
	if rc == nil {
		os.Exit(1)
	}
}
