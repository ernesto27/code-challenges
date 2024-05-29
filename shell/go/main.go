package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Current directory:", dir)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		for {
			<-sigs
			fmt.Print("\nccsh> ")
		}
	}()

	history := NewHistory()
	defer history.file.Close()

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
				break
			}
			continue
		}

		var output []byte

		if command[0] == "history" {
			for _, value := range history.currentCommands {
				output = append(output, []byte(fmt.Sprintf("%s\n", value))...)
			}

		} else {
			var cmd *exec.Cmd
			if strings.Contains(input, "|") {
				cmd = exec.Command("bash", "-c", input)
			} else {
				cmd = exec.Command(command[0], command[1:]...)
			}

			cmd.Dir = dir
			output, err = cmd.Output()

			if err != nil {
				fmt.Println("No such file or directory (os error 2)")
				continue
			}
		}

		history.currentCommands = append(history.currentCommands, input)
		_, err = history.file.WriteString(input + "\n")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Print(string(output))

	}
}

type History struct {
	file            *os.File
	currentCommands []string
}

func NewHistory() *History {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	homeDir := currentUser.HomeDir
	configFolder := ".config/ccsh"
	configFolderPath := filepath.Join(homeDir, configFolder)

	if _, err := os.Stat(configFolderPath); os.IsNotExist(err) {
		err := os.Mkdir(configFolderPath, 0700)
		if err != nil {
			panic(err)
		}
	}
	fileData := configFolderPath + "/" + "ccsh.txt"
	file, err := os.OpenFile(fileData, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	historyCommands := []string{}
	for scanner.Scan() {
		historyCommands = append(historyCommands, scanner.Text())
	}

	return &History{
		file:            file,
		currentCommands: historyCommands,
	}
}
