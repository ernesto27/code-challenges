package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateMasterPassword(t *testing.T) {
	vault := Vault{}
	p, err := vault.createMasterPassword("")
	if err != nil {
		t.Error("Error while creating master password")
	}

	if p == "" {
		t.Error("Error while creating master password")
	}

}

func TestCompareMasterPassword(t *testing.T) {
	vault := Vault{}
	hashedPassword := "$2a$10$ipt6MTq3qM6FLKPeTgKXp.4Btx1VXy6BnLrM9mfZdYIn1tulM0LzC"
	err := vault.compareMasterPassword(hashedPassword, "mypass")
	if err != nil {
		t.Error("Error while comparing master password")
	}

	err = vault.compareMasterPassword(hashedPassword, "wrongpass")
	if err == nil {
		t.Error("Error while comparing master password")
	}

	// invalid hash
	err = vault.compareMasterPassword("invalidhash", "mypass")
	if err == nil {
		t.Error("Error while comparing master password")
	}
}

func TestCreateVault(t *testing.T) {
	vault, err := newVault()
	if err != nil {
		t.Errorf("Error creating vault: %v", err)
	}
	name := "test"
	password := "testpass"

	filePath := filepath.Join(vault.configFolderPath, name)
	os.Remove(filePath)

	err = vault.createVault(name, password)
	if err != nil {
		t.Errorf("Error creating vault: %v", err)
	}

	// Vault file already created
	err = vault.createVault(name, password)
	if err == nil {
		t.Error("Should return error when vault already exists")
	}
}

func TestSignIn(t *testing.T) {
	vault, err := newVault()
	if err != nil {
		t.Errorf("Error creating vault: %v", err)
	}

	name := "signitTest"
	password := "1111"

	filePath := filepath.Join(vault.configFolderPath, name)
	os.Remove(filePath)

	err = vault.createVault(name, password)
	if err != nil {
		t.Errorf("Error creating vault: %v", err)
	}

	err = vault.signIn(name, password)
	if err != nil {
		t.Errorf("Error signing in: %v", err)
	}

	err = vault.signIn(name, "wrongpass")
	if err == nil {
		t.Error("Should return error when wrong password")
	}
}

func TestAddPassword(t *testing.T) {
	vault, err := newVault()
	if err != nil {
		t.Errorf("Error creating vault: %v", err)
	}

	name := "addPasswordTest"
	password := "1111"

	filePath := filepath.Join(vault.configFolderPath, name)
	os.Remove(filePath)

	err = vault.createVault(name, password)
	if err != nil {
		t.Errorf("Error creating vault: %v", err)
	}

	err = vault.signIn(name, password)
	if err != nil {
		t.Errorf("Error signing in: %v", err)
	}

	err = vault.addPassword("facebook", "myuser", "testpass")
	if err != nil {
		t.Errorf("Error adding password: %v", err)
	}
}

func TestGetPassword(t *testing.T) {
	vault, err := newVault()
	if err != nil {
		t.Errorf("Error creating vault: %v", err)
	}

	name := "getPasswordTest"
	password := "1111"

	filePath := filepath.Join(vault.configFolderPath, name)
	os.Remove(filePath)

	err = vault.createVault(name, password)
	if err != nil {
		t.Errorf("Error creating vault: %v", err)
	}

	err = vault.signIn(name, password)
	if err != nil {
		t.Errorf("Error signing in: %v", err)
	}

	err = vault.addPassword("facebook", "myuser", "testpass")
	if err != nil {
		t.Errorf("Error adding password: %v", err)
	}

	user, pass, err := vault.getPassword("facebook")
	if err != nil {
		t.Errorf("Error getting password: %v", err)
	}

	if user != "myuser" || pass != "testpass" {
		t.Error("Error getting password")
	}

	_, _, err = vault.getPassword("wrongpass")
	if err == nil {
		t.Error("Should return error when password not found")
	}
}
