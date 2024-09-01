package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/alexmullins/zip"
)

type ZipCracker struct {
	fileName     string
	file         *os.File
	passwordFile zip.ReadCloser
}

func NewZipCracker(file string) (*ZipCracker, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return &ZipCracker{file: f, fileName: file}, nil
}

func (z *ZipCracker) isValid() bool {
	signature := make([]byte, 4)
	_, err := z.file.Read(signature)
	if err != nil {
		return false
	}

	if signature[0] == 0x50 && signature[1] == 0x4B && signature[2] == 0x03 && signature[3] == 0x04 {
		return true
	}
	return false
}

func (z *ZipCracker) findPassword() {
	file, err := os.Open("realhuman_phill.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Read the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		password := scanner.Text()
		fmt.Println("Trying password: ", password)
		found, err := z.tryPassword(password)

		if err != nil {
			log.Fatal(err)
		}

		if found {
			fmt.Println("Password found: ", password)
			return
		}

	}
}

func (z *ZipCracker) tryPassword(password string) (bool, error) {
	r, err := zip.OpenReader(z.fileName)
	if err != nil {
		return false, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword(password)
		}
		rc, err := f.Open()
		if err == nil {
			rc.Close()
			return true, nil
		}
	}
	return false, nil
}

func main() {
	zipCracker, err := NewZipCracker("go.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer zipCracker.file.Close()

	if zipCracker.isValid() {
		fmt.Println("Valid ZIP file")
	} else {
		fmt.Println("Invalid ZIP file")
	}

	zipCracker.findPassword()

	v, err := zipCracker.tryPassword("test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v)
}
