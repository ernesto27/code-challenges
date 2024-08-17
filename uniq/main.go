package main

import (
	"fmt"
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

	return response
}

func main() {

	// scanner := bufio.NewScanner(os.Stdin)

	// standarInput := ""

	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	standarInput += fmt.Sprintf("%s", line)
	// }

	// fmt.Println(standarInput)

	text := `line1
line2
line2
line3
line1
line4
line2
line4
line5`

	uniq := Uniq{}
	fmt.Println(uniq.Count(text))

	// 	resp := uniq.Run(text)
	// 	fmt.Print(resp)
	// 	err := uniq.Save("output.txt")
	// 	if err != nil {
	// 		panic(err)
	// 	}

}
