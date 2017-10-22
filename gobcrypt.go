// gobcrypt.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
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

func main() {

	var count, pwLength int
	var thisPW, inFile, outFile string
	var testHash bool
	flag.IntVar(&count, "c", 1, "Specify the number of hashes to create")
	flag.StringVar(&thisPW, "s", "", "Hash the specified password only (1)")
	flag.BoolVar(&testHash, "t", false, "Validate the hash & pass")
	flag.IntVar(&pwLength, "l", 15, "Specify the length of password required")
	flag.StringVar(&inFile, "f", "", "Specify a file to read passwords from")
	flag.StringVar(&outFile, "o", "", "Specify a file to write out the pass/hash to")
	flag.Parse()

	if thisPW == "" {

		if inFile != "" {

			lines, err := readLines(inFile)
			
			if err != nil {
				log.Fatalf("readLines: %s", err)
			}
			var hashLines []string
			
			for _, line := range lines {
				
				password := line
				hash, _ := hashPassword(password)
				var matchLine string
				hashArray := []string{password, hash}
				hashLine := strings.Join(hashArray, "	")
				hashLines = append(hashLines, hashLine)

				if testHash != false {
					matchLine = matchPasswordAndHash(password, hash)
					hashLines = append(hashLines, matchLine)
				}
				
				if outFile != "" {
					if err := writeLines(hashLines, outFile); err != nil {
						log.Fatal("writeLine: %s", err)
					}
				} else {
					fmt.Println(hashLine)
					if testHash != false {
						fmt.Println(matchLine)
					}
				}
			}
		} else {

			var matchLine string
			var hashLines []string

			for count > 0 {

				password := randomString(pwLength)
				hash, _ := hashPassword(password) // TODO add error handling
				hashArray := []string{password, hash}
				hashLine := strings.Join(hashArray, "	")
				hashLines = append(hashLines, hashLine)

				if testHash != false {
					matchLine = matchPasswordAndHash(password, hash)
					hashLines = append(hashLines, matchLine)
				}
				if outFile == "" {
					fmt.Println(hashLine)
					if testHash != false {
						fmt.Println(matchLine)
					}
				}
				count -= 1
			}
			if outFile != "" {
				if err := writeLines(hashLines, outFile); err != nil {
					log.Fatal("writeLine: %s", err)
				}
			} 			
		}
	} else {

		if inFile != "" {
			log.Fatalf("Cannot run -s and -f together please select one")
		}

		var matchLine string
		var hashLines []string
		password := thisPW
		hash, _ := hashPassword(password)
		hashArray := []string{password, hash}
		hashLine := strings.Join(hashArray, "	")
		hashLines = append(hashLines, hashLine)

		if testHash != false {
			matchLine = matchPasswordAndHash(password, hash)
			hashLines = append(hashLines, matchLine)
		}
		if outFile != "" {
			if err := writeLines(hashLines, outFile); err != nil {
				log.Fatal("writeLine: %s", err)
			}
		} else {
			fmt.Println(hashLine)
			if testHash != false {
				fmt.Println(matchLine)
			}
		}
	}
}
