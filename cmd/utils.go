package main

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// errMalformedLine is returned when a pass/hash line does not contain a
// "<password> <hash>" pair.
var errMalformedLine = errors.New("malformed line: expected \"<password> <hash>\"")

const randomChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// randomString returns a cryptographically secure random string of strlen
// characters drawn from randomChars.
func randomString(strlen int) (string, error) {
	out := make([]byte, strlen)
	max := big.NewInt(int64(len(randomChars)))
	for i := range out {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		out[i] = randomChars[n.Int64()]
	}
	return string(out), nil
}

// hashPassword creates a bcrypt hash of password using the given cost.
func hashPassword(password string, cost int) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(b), err
}

// splitPassHash splits a "<password> <hash>" line on its last space. Bcrypt
// hashes never contain spaces, but passwords may, so the hash is always the
// final field.
func splitPassHash(line string) (password, hash string, err error) {
	i := strings.LastIndex(line, " ")
	if i < 0 {
		return "", "", errMalformedLine
	}
	return line[:i], line[i+1:], nil
}

// checkHashAndPassword reports whether the password on a "<password> <hash>"
// line matches its hash.
func checkHashAndPassword(hashLine string) error {
	password, hash, err := splitPassHash(hashLine)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// readLines reads a file into memory and returns its non-empty lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

// writeLines writes lines to path. The file is created with mode 0600 because
// it contains plaintext passwords.
func writeLines(lines []string, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	if err := w.Flush(); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

// returnOutput writes lines to outFile, or to stdout when outFile is empty.
func returnOutput(lines []string, outFile string) error {
	if outFile != "" {
		return writeLines(lines, outFile)
	}
	for _, line := range lines {
		fmt.Println(line)
	}
	return nil
}

// generateHashForPassword returns a single "<password> <hash>" line.
func generateHashForPassword(password string, cost int) ([]string, error) {
	hash, err := hashPassword(password, cost)
	if err != nil {
		return nil, err
	}
	return []string{password + " " + hash}, nil
}

// hashPasswordsFromFile hashes each password. Passwords longer than 72 bytes
// are counted in skipped and left out; any other error aborts.
func hashPasswordsFromFile(passwords []string, cost int) (lines []string, skipped int, err error) {
	for _, password := range passwords {
		hashed, err := generateHashForPassword(password, cost)
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			skipped++
			continue
		}
		if err != nil {
			return nil, skipped, err
		}
		lines = append(lines, hashed...)
	}
	return lines, skipped, nil
}

// matchPasswordAndHash validates a "<password> <hash>" line and formats the
// result. The returned error is non-nil if the line is malformed or does not
// match.
func matchPasswordAndHash(hashLine string) (string, error) {
	password, hash, err := splitPassHash(hashLine)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Sprintf("MATCH: FAIL, password: %s, hash: %s", password, hash), err
	}
	return fmt.Sprintf("MATCH: PASS, password: %s, hash: %s", password, hash), nil
}
