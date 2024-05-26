package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Current directory:", dir)

	for {
		fmt.Print("ccsh> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		command := strings.Split(input, " ")

		if command[0] == "exit" {
			break
		}

		if len(command) == 2 && command[0] == "cd" {
			err := os.Chdir(command[1])
			if err != nil {
				fmt.Println("No such file or directory (os error 2)")
			}
			dir, err = os.Getwd()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			continue

		}

		cmd := exec.Command(command[0], command[1:]...)
		cmd.Dir = dir
		output, err := cmd.Output()

		if err != nil {
			fmt.Println("No such file or directory (os error 2)")
			continue
		}

		fmt.Print(string(output))
	}
}
