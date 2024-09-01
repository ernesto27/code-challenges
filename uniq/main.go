package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Uniq struct {
	response string
}

func (uniq *Uniq) Run(text string) string {
	m := make(map[string]bool)
	lines := strings.Split(text, "\n")

	var response string
	for _, line := range lines {
		if m[line] || line == "" {
			continue
		}
		m[line] = true
		response += fmt.Sprintf("%s\n", line)
	}

	response = strings.TrimSuffix(response, "\n")
	uniq.response = response
	return response
}

func (uniq *Uniq) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(uniq.response)
	if err != nil {
		return err
	}

	return nil
}

func (uniq *Uniq) Count(text string) string {
	lines := strings.Split(text, "\n")

	type item struct {
		index int
		count int
		value string
	}

	itemsSlice := []item{}
	m := make(map[string]item)
	indexCount := 0
	for _, line := range lines {
		if line == "" {
			continue
		}

		if _, ok := m[line]; ok {
			currItem := item{
				index: m[line].index,
				count: m[line].count + 1,
				value: line,
			}
			m[line] = currItem
			itemsSlice[m[line].index] = currItem
		} else {
			currItem := item{
				index: indexCount,
				count: 1,
				value: line,
			}
			indexCount++

			m[line] = currItem
			itemsSlice = append(itemsSlice, currItem)
		}
	}

	var response string
	for _, item := range itemsSlice {
		response += fmt.Sprintf("%d %s\n", item.count, item.value)
	}

	response = strings.TrimSuffix(response, "\n")
	return response
}

func (uniq *Uniq) getDuplicated(text string) string {
	m := make(map[string]int)
	lines := strings.Split(text, "\n")

	var response string
	for _, line := range lines {
		if _, ok := m[line]; ok {
			m[line]++

			if m[line] == 2 {
				response += fmt.Sprintf("%s\n", line)
			}
			continue
		}
		m[line] = 1
	}

	response = strings.TrimSuffix(response, "\n")
	uniq.response = response
	return response
}

func main() {
	outputPathFile := flag.String("o", "", "output path file")
	count := flag.Bool("c", false, "show count")
	duplicated := flag.Bool("d", false, "show duplicated")
	flag.Parse()

	uniq := Uniq{}
	var pathFile string
	if len(os.Args) >= 2 {
		pathFile = os.Args[len(os.Args)-1]
	}

	if pathFile != "" {
		fileContent, err := os.ReadFile(pathFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		text := string(fileContent)
		resp := uniq.Run(text)

		if *count {
			resp = uniq.Count(text)
		}

		if *duplicated {
			resp = uniq.getDuplicated(text)
		}

		fmt.Println(resp)
	} else {
		stdinContent, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		text := string(stdinContent)
		resp := uniq.Run(text)

		if *outputPathFile != "" {
			err := uniq.Save(*outputPathFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println(resp)
		}
	}

}
