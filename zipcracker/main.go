package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type ZipCracker struct {
	fileName     string
	fileZip      *os.File
	filePassword *os.File
	minLen       int
	maxLen       int
}

func NewZipCracker(fileNameZip string, filePasswordName string, minLen, maxLen int) (*ZipCracker, error) {
	fileZip, err := os.Open(fileNameZip)
	if err != nil {
		return nil, err
	}

	filePassword, err := os.Open(filePasswordName)
	if err != nil {
		return nil, err
	}

	return &ZipCracker{
		fileName:     fileNameZip,
		fileZip:      fileZip,
		filePassword: filePassword,
		minLen:       minLen,
		maxLen:       maxLen,
	}, nil
}

func (z *ZipCracker) isValid() bool {
	signature := make([]byte, 4)
	_, err := z.fileZip.Read(signature)
	if err != nil {
		return false
	}

	if signature[0] == 0x50 && signature[1] == 0x4B && signature[2] == 0x03 && signature[3] == 0x04 {
		return true
	}
	return false
}

func (z *ZipCracker) findPassword() (string, error) {
	scanner := bufio.NewScanner(z.filePassword)
	var wg sync.WaitGroup
	var foundPassword string
	semaphore := make(chan struct{}, 5)

	for scanner.Scan() {
		if foundPassword != "" {
			return foundPassword, nil
		}

		passwordsToTry := []string{}
		originalPassword := scanner.Text()
		passwordsToTry = append(passwordsToTry, originalPassword)

		if z.minLen >= 1 && z.maxLen >= len(originalPassword)-1 {
			combinePasswords := z.generateCombinations(scanner.Text())
			passwordsToTry = append(passwordsToTry, combinePasswords...)
		}

		for _, password := range passwordsToTry {
			fmt.Println("Trying password: ", password)
			semaphore <- struct{}{}
			wg.Add(1)

			go func(password string) {
				defer wg.Done()
				defer func() { <-semaphore }()
				valid := z.tryPassword(password)
				if valid {
					foundPassword = password
				}
			}(password)
		}

	}

	wg.Wait()

	if foundPassword != "" {
		return foundPassword, nil
	}

	return "", fmt.Errorf("password not found")

}

func (z *ZipCracker) tryPassword(password string) bool {
	cmd := exec.Command("unzip", "-P", password, "-t", z.fileName)
	_, err := cmd.CombinedOutput()

	return err == nil
}

func (z *ZipCracker) generateCombinations(password string) []string {
	chars := strings.Split(password, "")
	n := len(chars)
	combinations := []string{}

	var generate func([]string, int)
	generate = func(current []string, position int) {
		if position == n {
			combinations = append(combinations, strings.Join(current, ""))
			return
		}

		for i := position; i < n; i++ {
			current[position], current[i] = current[i], current[position]
			generate(current, position+1)
			current[position], current[i] = current[i], current[position]
		}
	}

	generate(chars, 0)
	return combinations
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <zip_file_name>")
		os.Exit(1)
	}

	minLen := flag.Int("min", 0, "Minimum length of combinations")
	maxLen := flag.Int("max", 0, "Maximum length of combinations")
	flag.Parse()

	if *minLen > *maxLen {
		fmt.Println("Error: min length cannot be greater than max length")
		flag.PrintDefaults()
		return
	}

	fmt.Println("Min length: ", *minLen)
	fmt.Println("Max length: ", *maxLen)
	fmt.Println("File name: ", os.Args)

	zipFileName := os.Args[len(os.Args)-1]
	zipCracker, err := NewZipCracker(zipFileName, "passwords.txt", *minLen, *maxLen)
	if err != nil {
		log.Fatal(err)
	}
	defer zipCracker.fileZip.Close()
	defer zipCracker.filePassword.Close()

	if zipCracker.isValid() {
		fmt.Println("Valid ZIP file")
	} else {
		fmt.Println("Invalid ZIP file")
	}

	resp, err := zipCracker.findPassword()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Password found: ", resp)

}
