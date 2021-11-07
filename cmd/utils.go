package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var r *rand.Rand // Rand for this package.

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// This function will generate a random string based on the random value
//  set from the init function ant the chars defined below
func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := ""
	for i := 0; i < strlen; i++ {
		index := r.Intn(len(chars))
		result += chars[index : index+1]
	}
	return result
}

/* This function will create a hashed password using bcrypt. It requires a
*  String value to be passed as an argument and returns the hash as a String
 */
func hashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

/* This function compares the plain text password with the hash to ensure it is
*  valid. It requires both the plain text and hash pass String to be passed to it
 */
func checkHashAndPassword(hashLine string) error {
	items := strings.Split(hashLine, " ")
	hash, password := items[1], items[0]
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

// hanldes returning the output to user via file or stdout
func returnOutput(hashLines []string, outFile string) error {
	if outFile != "" {
		if err := writeLines(hashLines, outFile); err != nil {
			return err
		}
	} else {
		for _, line := range hashLines {
			fmt.Println(line)
		}
	}
	return nil
}

// generates a hash from the given password string
func generateHashForPassword(password string, cost int) ([]string, error) {
	var hashLines []string

	hash, err := hashPassword(password, cost)
	if err != nil {
		return nil, err
	}

	hashLine := strings.Join([]string{password, hash}, " ")
	hashLines = append(hashLines, hashLine)

	return hashLines, err
}

// Processes the returned Bool value from checkHashAndPassword
//  and formats it for use in hashLines
func matchPasswordAndHash(hashLine string) (result string, err error) {
	items := strings.Split(hashLine, " ")
	hash, password := items[1], items[0]
	err = checkHashAndPassword(hashLine)
	if err != nil {
		result = fmt.Sprintf("MATCH: FAIL, password: %s, hash: %s", password, hash)
		return result, err
	}
	result = fmt.Sprintf("MATCH: PASS, password: %s, hash: %s", password, hash)
	return result, err
}
