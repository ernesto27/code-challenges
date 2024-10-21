package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type Head struct {
	files []*os.File
}

func newHead(files []string) (*Head, error) {
	if len(files) == 0 {
		return &Head{}, nil
	}

	var fileDesc []*os.File

	for _, file := range files {
		file, err := os.Open(file)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return nil, err
		}
		fileDesc = append(fileDesc, file)
	}

	return &Head{
		files: fileDesc,
	}, nil
}

func (h *Head) Run(lines int, bytes int) (string, error) {
	var result string
	for _, file := range h.files {
		if len(h.files) > 1 {
			result += "==> " + file.Name() + " <==\n"
		}

		reader := bufio.NewReader(file)
		lineCount := 0

		if bytes > 0 {
			buf := make([]byte, bytes)
			n, err := reader.Read(buf)
			if err != nil && err != io.EOF {
				return "", err
			}
			return string(buf[:n]), nil

		} else {
			for lineCount < lines {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						result += line
						break
					}
					return "", err
				}
				result += line
				lineCount++
			}
		}
	}

	return result, nil
}

func (h *Head) stdInputPrint() {
	count := 0
	for count < 10 {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		fmt.Println(input)
		count++
	}
}

func (h *Head) Close() {
	for _, file := range h.files {
		file.Close()
	}
}

func main() {
	limit := flag.Int("n", 10, "Number of lines to print")
	bytes := flag.Int("c", 0, "Number of bytes to print")
	flag.Parse()

	fileArgs := flag.Args()
	if len(fileArgs) == 0 {
		head, _ := newHead([]string{})
		head.stdInputPrint()
		return
	}

	head, err := newHead(fileArgs)
	if err != nil {
		fmt.Println("Error creating head:", err)
		os.Exit(1)
	}
	defer head.Close()

	resp, err := head.Run(*limit, *bytes)
	if err != nil {
		fmt.Println("Error running head:", err)
		os.Exit(1)
	}

	fmt.Print(resp)

}
