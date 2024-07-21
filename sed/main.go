package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Sed struct {
	input          string
	pattern        string
	replace        string
	fromRangePrint int
	toRangePrint   int
}

func newSed(cmd, input string) (*Sed, error) {
	values := strings.Split(cmd, "/")
	var patternVal, replacementVal string
	var fromRange, toRange int

	if len(values) >= 3 {
		patternVal = values[1]
		replacementVal = values[2]
	} else {
		cmd = cmd[:len(cmd)-1]
		values = strings.Split(cmd, ",")

		if len(values) == 2 {
			f, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, err
			}
			t, err := strconv.Atoi(values[1])
			if err != nil {
				return nil, err
			}
			fromRange = f
			toRange = t
		} else {
			patternVal = values[0]
		}
	}

	return &Sed{
		pattern:        patternVal,
		replace:        replacementVal,
		input:          input,
		fromRangePrint: fromRange,
		toRangePrint:   toRange,
	}, nil
}

func (s *Sed) doReplace() string {
	pattern := regexp.MustCompile(s.pattern)
	output := pattern.ReplaceAllString(s.input, s.replace)
	return output
}

func (s *Sed) rangeLines() string {
	var resp string
	lines := strings.Split(s.input, "\n")

	for idx, line := range lines {
		idx = idx + 1
		if idx >= s.fromRangePrint && idx <= s.toRangePrint {
			resp += line + "\n"
		}
	}

	resp = strings.TrimSuffix(resp, "\n")
	return resp
}

func (s *Sed) getLinesByPattern() string {
	pattern := regexp.MustCompile(s.pattern)
	lines := strings.Split(s.input, "\n")
	var resp string

	for _, line := range lines {
		if pattern.MatchString(line) {
			resp += line + "\n"
		}
	}

	resp = strings.TrimSuffix(resp, "\n")

	return resp
}

func (s *Sed) doubleSpacing() string {
	lines := strings.Split(s.input, "\n")
	var resp string

	for idx, line := range lines {
		if idx == len(lines)-1 {
			resp += line
		} else {
			resp += line + "\n\n"
		}
	}

	return resp
}

func (s *Sed) removeTrailingSpaces() string {
	resp := strings.TrimSpace(s.input)
	return resp
}

func main() {

	n := flag.String("n", "", "Print the line number")
	flag.Parse()

	if *n == "" {

		pattern := os.Args[1]
		file := os.Args[2]

		f, err := os.Open(file)
		if err != nil {
			fmt.Println("Error opening file: ", err)
			os.Exit(1)
		}
		defer f.Close()

		fileContent, err := io.ReadAll(f)
		if err != nil {
			fmt.Println("Error reading file: ", err)
			os.Exit(1)
		}

		sed, err := newSed(pattern, string(fileContent))
		if err != nil {
			fmt.Println("Error creating Sed object: ", err)
			os.Exit(1)
		}

		fmt.Println(sed.doReplace())
	} else {
		reader := os.Stdin
		stdContent, err := io.ReadAll(reader)
		if err != nil {
			fmt.Println("Error reading input: ", err)
			os.Exit(1)
		}

		sed, err := newSed(*n, string(stdContent))
		if err != nil {
			fmt.Println("Error creating Sed object: ", err)
			os.Exit(1)
		}

		resp := sed.rangeLines()
		fmt.Println(resp)
	}

}
