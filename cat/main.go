package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Cat struct {
	files []*os.File
}

func NewCat(fileName ...string) *Cat {
	return &Cat{}
}

func (c *Cat) WithFiles(fileName ...string) (*Cat, error) {
	files := []*os.File{}

	for _, file := range fileName {
		if strings.HasPrefix(file, "-") {
			continue
		}
		file, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	c.files = append(c.files, files...)
	return c, nil
}

func (c *Cat) WithStdin(stdin *os.File) (*Cat, error) {
	c.files = append(c.files, os.Stdin)
	return c, nil
}

func (c *Cat) GetContentFile(showLineNumbers, ignoreEmptyLines bool) (string, error) {
	var resp string
	for _, file := range c.files {
		var buffer bytes.Buffer
		reader := bufio.NewReader(file)

		if showLineNumbers || ignoreEmptyLines {
			scanner := bufio.NewScanner(reader)

			lineNumber := 1
			for scanner.Scan() {
				lineNumberStr := fmt.Sprintf("%6d", lineNumber)
				if ignoreEmptyLines && strings.TrimSpace(scanner.Text()) == "" {
					resp += "\n"
				} else {
					resp += fmt.Sprintf("%s %s\n", lineNumberStr, scanner.Text())
					lineNumber++
				}
			}
		} else {
			_, err := buffer.ReadFrom(reader)
			if err != nil {
				return "", err
			}
			resp += buffer.String()
		}
	}

	return resp, nil
}

func (c *Cat) Close() {
	for _, file := range c.files {
		file.Close()
	}
}

func main() {
	showLineNumbers := flag.Bool("n", false, "show line numbers")
	ignoreEmptyLines := flag.Bool("b", false, "ignore empty lines")
	flag.Parse()

	var content string
	cat := &Cat{}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var err error
		cat, err = NewCat().WithStdin(os.Stdin)
		if err != nil {
			fmt.Println("Error reading stdin:", err)
			os.Exit(1)
		}
	} else {
		if flag.NArg() < 1 {
			fmt.Println("Please provide a filename as an argument")
			os.Exit(1)
		}
		var err error
		cat, err = NewCat().WithFiles(flag.Args()...)
		if err != nil {
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
	}

	defer cat.Close()

	content, err := cat.GetContentFile(*showLineNumbers, *ignoreEmptyLines)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	fmt.Print(content)
}
