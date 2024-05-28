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

	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	homeDir := currentUser.HomeDir
	configFolder := ".config/ccsh"
	configFolderPath := filepath.Join(homeDir, configFolder)
	fmt.Println(configFolderPath)

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
	fmt.Println(file)

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
