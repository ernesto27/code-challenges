package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	param := os.Args[1]
	// fmt.Println(param)
	file := os.Args[2]

	switch param {
	case "-c":

		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%d %s\n", fileInfo.Size(), fileInfo.Name())

	case "-l":
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		lineCount := bytes.Count(content, []byte("\n"))
		fmt.Printf("%d %s\n", lineCount, file)

	}

}
