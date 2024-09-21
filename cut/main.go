package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Cut struct {
	fileDesc   *os.File
	recordList map[int][]string
	delimiter  string
}

func NewCut(filename string) *Cut {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	return &Cut{fileDesc: file}
}

func NewCutWithStdin() *Cut {
	return &Cut{fileDesc: os.Stdin}
}

func (c *Cut) Read(index string, delimiter string) error {
	if index == "" {
		return fmt.Errorf("index out of range")
	}
	reader := csv.NewReader(c.fileDesc)
	c.delimiter = delimiter

	if delimiter != "\t" {
		newDelimenter := []rune(delimiter)
		reader.Comma = newDelimenter[0]
	} else {
		reader.Comma = '\t'
	}

	indexList := []int{}
	indexToPrint := strings.Split(index, ",")
	for _, i := range indexToPrint {
		indexVal, err := strconv.Atoi(i)
		if err != nil {
			return fmt.Errorf("not valid index")
		}

		if indexVal <= 0 {
			return fmt.Errorf("not valid index")
		}

		indexList = append(indexList, indexVal)
	}

	m := make(map[int][]string)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		for _, currIndex := range indexList {
			if len(record) < currIndex {
				return fmt.Errorf("index out of range 2")
			}
			m[currIndex-1] = append(m[currIndex-1], record[currIndex-1])
		}

	}
	c.recordList = m
	return nil
}

func (c *Cut) Print() {
	var keys []int
	for k := range c.recordList {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		fmt.Printf("%s%s", c.recordList[k][0], c.delimiter)
	}
	fmt.Println()

	for i := 1; i < len(c.recordList[keys[0]]); i++ {
		for _, k := range keys {
			fmt.Printf("%s%s", c.recordList[k][i], c.delimiter)
		}
		fmt.Println()
	}
}

func (c *Cut) Close() {
	c.fileDesc.Close()
}

func main() {
	fieldIndex := flag.String("f", "", "field number to display (1-based index)")
	delimiter := flag.String("d", "\t", "delimiter")

	flag.Parse()

	fileName := flag.Args()

	stat, _ := os.Stdin.Stat()
	var cut *Cut
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		cut = NewCutWithStdin()
		if cut == nil {
			return
		}

	} else {

		if len(fileName) == 0 {
			fmt.Println("Please provide a filename as an argument")
			return
		}

		cut = NewCut(fileName[len(fileName)-1])
		if cut == nil {
			return
		}
	}

	defer cut.Close()

	err := cut.Read(*fieldIndex, *delimiter)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	cut.Print()
}
