package main

import (
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
	} else if len(values) == 2 {
		cmd = cmd[:len(cmd)-1]
		values = strings.Split(cmd, ",")

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

func main() {
	//cmd := "s/\"/44/g"
	cmd := "/roads/p"

	text := "\"Your heart is the size of an ocean. Go find yourself in is hidden depths size.\""
	sed, err := newSed(cmd, text)
	if err != nil {
		panic(err)
	}

	sed.doReplace()

	// values := strings.Split(cmd, "/")

	// patternVal := values[1]
	// replacementVal := values[2]

	// text := "\"Your heart is the size of an ocean. Go find yourself in is hidden depths size.\""

	// pattern := regexp.MustCompile(patternVal)

	// // Replace all occurrences of "this" with "that"
	// replacedText := pattern.ReplaceAllString(text, replacementVal)

	// fmt.Println(replacedText)

}
