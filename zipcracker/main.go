package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

type ZipCracker struct {
	fileName     string
	fileZip      *os.File
	filePassword *os.File
}

func NewZipCracker(fileNameZip string, filePasswordName string) (*ZipCracker, error) {
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

		semaphore <- struct{}{}
		password := scanner.Text()
		fmt.Println("Trying password: ", password)

		wg.Add(1)

		go func(password string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			valid := z.tryPassword(password)
			if valid {
				// fmt.Println("Password found: ", password)
				foundPassword = password
				return
			}
		}(password)

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

func main() {

	zipCracker, err := NewZipCracker("file.zip", "realhuman_phill.txt")
	// zipCracker, err := NewZipCracker("file.zip", "password.txt")
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

	// valid := zipCracker.tryPassword("tes1t")
	// if !valid {
	// 	fmt.Println("Password is not valid")

	// } else {
	// 	fmt.Println(valid)

	// }
}
