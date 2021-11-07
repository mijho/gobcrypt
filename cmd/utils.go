package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
	"strings"
	"time"
	"strconv"
)

var r *rand.Rand // Rand for this package.

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

/* This function will generate a random string based on the random value
*  set from the init function ant the chars defined below
 */
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
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

/* This function compares the plain text password with the hash to ensure it is
*  valid. It requires both the plain text and hash pass String to be passed to it
 */
func checkHashAndPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

func writeSourcesToFile(sources []string, outFile string) error {
	file, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	datawriter := bufio.NewWriter(file)

	for _, line := range sources {
		_, _ = datawriter.WriteString(line + "\n")
	}

	datawriter.Flush()
	file.Close()

	return err
}

/* Processes the returned Bool value from checkHashAndPassword
*  and formats it for use in hashLines
*/
func matchPasswordAndHash(password, hash string) string {
	match := checkHashAndPassword(password, hash)
	matchString := strconv.FormatBool(match)
	matchArray := []string{"Match: ", matchString}
	matchLine := strings.Join(matchArray, " ")
	return string(matchLine)
}