package main

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func Test_generateHashForPassword(t *testing.T) {
	password := "superSecretPassword123"
	hashLine, _ := generateHashForPassword(password, 14)
	if err := checkHashAndPassword(hashLine[0]); err != nil {
		t.Error("Hash and Password do not match")
	}
}

func Test_generateHashForRandomPassword(t *testing.T) {
	password := randomString(20)
	hashLine, _ := generateHashForPassword(password, 14)
	if err := checkHashAndPassword(hashLine[0]); err != nil {
		t.Error("Hash and Password do not match")
	}
}

func Test_InvalidHashAndPassword(t *testing.T) {
	password := "password"
	hash := "$2a$14$9/5dnFqcgsFRWYEOWi3iyiyiaPOBX6yEJA4vNZIK78cUIOuryK3vMPYfSWs"
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		t.Errorf("Password and Hash should not match: %v", err)
	}
}

func Test_RandomPasswordLength(t *testing.T) {
	password := randomString(20)
	if len(password) != 20 {
		t.Error("Generated Password not correct length")
	}
}
