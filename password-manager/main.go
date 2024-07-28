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

	file, err := v.getFileVault(fileData)
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

	// TODO CHECK NAME ALREADY PRESENT

	file.Write([]byte(name + ":" + username + ":" + encryptedPassword + "\n"))

	return nil
}

func (v *Vault) getPassword(name string) (string, error) {
	fileData := filepath.Join(v.configFolderPath, v.name)
	file, err := v.getFileVault(fileData)
	if err != nil {
		return "", err
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
			return "", errors.New("invalid vault file")
		}

		if split[0] == name {
			decryptedPassword, err := v.decrypt(split[2], v.secretKey)
			if err != nil {
				return "", err
			}
			return decryptedPassword, nil
		}
	}

	return "", errors.New("password not found")

}

func (v *Vault) getFileVault(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
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

	v, err := newVault()
	if err != nil {
		panic(err)
	}

	// v.createVault("personal", "1111")

	err = v.signIn("personal", "1111")
	if err != nil {
		panic(err)
	}

	s, err := v.getPassword("googlde")
	if err != nil {
		panic(err)
	}
	fmt.Println(s)

	// v.addPassword("facebook", "user", "password")
	// v.addPassword("google", "user", "1111")
}
