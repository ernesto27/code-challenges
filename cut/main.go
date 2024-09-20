package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Cut struct {
	fileDesc *os.File
	records  []string
}

func NewCut(filename string) *Cut {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	return &Cut{fileDesc: file}
}

func (c *Cut) Read(index string, delimiter string) ([]string, error) {
	if index == "" {
		return nil, fmt.Errorf("index out of range")
	}
	reader := csv.NewReader(c.fileDesc)

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
			return nil, fmt.Errorf("not valid index")
		}

		if indexVal <= 0 {
			return nil, fmt.Errorf("not valid index")
		}

		indexList = append(indexList, indexVal)
	}

	resp := []string{}
	resp2 := make([][]string, len(indexList))
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		for _, currIndex := range indexList {
			if len(record) <= currIndex {
				return resp, fmt.Errorf("index out of range 2")
			}
			fmt.Println("loop:", currIndex-1)
			resp = append(resp, record[currIndex-1])
			resp2[currIndex-1] = append(resp2[currIndex-1], record[currIndex-1])
		}

	}
	c.records = resp

	fmt.Println("resp2:", resp2)
	return resp, nil
}

func (c *Cut) Print() {
	for _, record := range c.records {
		fmt.Println(record)
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
	if len(fileName) == 0 {
		fmt.Println("Please provide a filename as an argument")
		return
	}

	cut := NewCut(fileName[len(fileName)-1])
	if cut == nil {
		return
	}
	defer cut.Close()

	_, err := cut.Read(*fieldIndex, *delimiter)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	cut.Print()

}
