package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

type Vault struct {
	name             string
	secretKey        []byte
	configFolderPath string
}

func newVault() (*Vault, error) {
	v := Vault{}
	configFolderPath, err := v.getConfigFolder()
	if err != nil {
		return &v, err
	}

	v.configFolderPath = configFolderPath
	return &v, nil
}

func (v *Vault) createMasterPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func (v *Vault) compareMasterPassword(hashPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) createVault(name string, password string) error {
	fileData := filepath.Join(v.configFolderPath, name)

	currentFile, _ := os.Stat(fileData)
	if currentFile != nil {
		return errors.New("Vault already exists")
	}

	file, err := v.createFileVault(fileData)
	if err != nil {
		return err
	}

	defer file.Close()

	hashPassword, err := v.createMasterPassword(password)
	if err != nil {
		return err
	}

	file.Write([]byte(name + ":" + hashPassword + "\n"))
	return nil
}

func (v *Vault) signIn(name string, password string) error {
	fileData := filepath.Join(v.configFolderPath, name)
	file, err := v.getFileVault(fileData)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println(file)

	scanner := bufio.NewScanner(file)
	var namePassword string
	if scanner.Scan() {
		namePassword = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	nameHashPassword := namePassword[len(name)+1:]

	err = v.compareMasterPassword(nameHashPassword, password)
	if err != nil {
		return err
	}

	secret := sha256.Sum256([]byte(password))
	v.name = name
	v.secretKey = secret[:]

	return nil
}

func (v *Vault) addPassword(name, username, password string) error {
	fileData := filepath.Join(v.configFolderPath, v.name)
	file, err := v.getFileVault(fileData)
	if err != nil {
		return err
	}
	defer file.Close()

	encryptedPassword, err := v.encrypt([]byte(password), []byte(v.secretKey))
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		split := strings.Split(line, ":")

		if len(split) != 3 {
			return errors.New("invalid vault file")
		}

		if split[0] == name {
			return errors.New("name already exists")
		}
	}

	file.Write([]byte(name + ":" + username + ":" + encryptedPassword + "\n"))

	return nil
}

func (v *Vault) getPassword(name string) (string, string, error) {
	fileData := filepath.Join(v.configFolderPath, v.name)
	file, err := v.createFileVault(fileData)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		split := strings.Split(line, ":")

		if len(split) != 3 {
			return "", "", errors.New("invalid vault file")
		}

		if split[0] == name {
			decryptedPassword, err := v.decrypt(split[2], v.secretKey)
			if err != nil {
				return "", "", err
			}
			return split[1], decryptedPassword, nil
		}
	}

	return "", "", errors.New("password not found")

}

func (v *Vault) createFileVault(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (v *Vault) getFileVault(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (v *Vault) getConfigFolder() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	homeDir := currentUser.HomeDir
	configFolder := ".config/password-manager"
	configFolderPath := filepath.Join(homeDir, configFolder)

	if _, err := os.Stat(configFolderPath); os.IsNotExist(err) {
		err := os.Mkdir(configFolderPath, 0700)
		if err != nil {
			return "", err
		}
	}

	return configFolderPath, nil
}

func (v *Vault) encrypt(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

func (v *Vault) decrypt(ciphertext string, key []byte) (string, error) {
	data, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], string(data[nonceSize:])
	plaintext, err := gcm.Open(nil, nonce, []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// save vault files on HOME dir
// save username - password pairs in a file enctrypted

func main() {
	vault, err := newVault()
	if err != nil {
		panic(err)
	}

	fmt.Println("Welcome to CC Password Manager")
	fmt.Println("What would you like to do?")

	var userInput string
	for {
		fmt.Println("1. Create Vault")
		fmt.Println("2. Sign In")
		fmt.Println("3. Add Password")
		fmt.Println("4. Get Password")
		fmt.Println("Quit (q)")

		fmt.Scanln(&userInput)

		switch userInput {
		case "1":
			input, err := getInputValues(false)
			if err != nil {
				continue
			}

			err = vault.createVault(input.name, input.password)
			if err != nil {
				fmt.Println("Error creating vault: ", err)
			}

		case "2":
			input, err := getInputValues(false)
			if err != nil {
				continue
			}

			err = vault.signIn(input.name, input.password)
			if err != nil {
				fmt.Println("Error signing in: ")
				continue
			}

			fmt.Println("Signed in successfully")

		case "3":
			input, err := getInputValues(true)
			if err != nil {
				continue
			}

			err = vault.addPassword(input.name, input.username, input.password)
			if err != nil {
				fmt.Println("Error adding password: ", err)
			}

		case "4":
			fmt.Println("Enter the name of the password you want to retrieve")
			var name string
			fmt.Scanln(&name)

			username, password, err := vault.getPassword(name)
			if err != nil {
				fmt.Println("Error getting password: ", err)
				continue
			}

			fmt.Println("Username: ", username)
			fmt.Println("Password: ", password)
			fmt.Println()

		case "q":
			fmt.Println("Goodbye!")
			os.Exit(0)
		}
	}

}

type input struct {
	name     string
	username string
	password string
}

func getInputValues(username bool) (input, error) {
	input := input{}

	fmt.Print("Enter the name: ")
	var name string
	fmt.Scanln(&name)

	if username {
		fmt.Print("Enter the username: ")
		var username string
		fmt.Scanln(&username)
		input.username = username
	}

	fmt.Print("Enter the password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading password: ")
		return input, err
	}
	password := string(passwordBytes)
	fmt.Println()

	if name == "" || password == "" {
		fmt.Println("Name and password cannot be empty")
		return input, err
	}

	input.name = name
	input.password = password
	return input, nil
}
