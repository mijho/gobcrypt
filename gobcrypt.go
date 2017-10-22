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
func RandomString2(strlen int) string {
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
func HashPassword(password string) (string, error) {
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

//func hashLines(lines []string) {
//	var hashlines []string
//	for _, line := range lines {
//		password := line
//		hash, _ := HashPassword(password)
//		hashlines := append(hashlines, password, "	", hash)
//	}
//	return hashlines
//}

func matchPasswordAndHash(password, hash string) string {
	match := checkHashAndPassword(password, hash)
	matchString := strconv.FormatBool(match)
	matchArray := []string{"Match: ", matchString}
	matchLine := strings.Join(matchArray, " ")
	return string(matchLine)
}



func main() {

	var count, pwlength int
	var thispw, infile, outfile string
	var testhash bool
	flag.IntVar(&count, "c", 1, "Specify the number of hashes to create")
	flag.StringVar(&thispw, "s", "", "Hash the specified password only (1)")
	flag.BoolVar(&testhash, "t", false, "Validate the hash & pass")
	flag.IntVar(&pwlength, "l", 15, "Specify the length of password required")
	flag.StringVar(&infile, "f", "", "Specify a file to read passwords from")
	flag.StringVar(&outfile, "o", "", "Specify a file to write out the pass/hash to")
	flag.Parse()

	if thispw == "" {
		for count > 0 {

			if infile != "" {
				lines, err := readLines(infile)
				if err != nil {
					log.Fatalf("readLines: %s", err)
				}
				// fmt.Println(lines)
				var hashlines []string
				for _, line := range lines {
					password := line
					hash, _ := HashPassword(password)
					matchLine := "hats"

					hasharray := []string{password, hash}
					hashline := strings.Join(hasharray, "	")
					hashlines = append(hashlines, hashline)

					// fmt.Println(password, "	", hash)

					if testhash != false {
						matchLine = matchPasswordAndHash(password, hash)
						hashlines = append(hashlines, matchLine)
					}

					//fmt.Println(hashline)
					
					if outfile != "" {
						if err := writeLines(hashlines, outfile); err != nil {
							log.Fatal("writeLine: %s", err)
						}
					} else {
						fmt.Println(hashline)
						fmt.Println(matchLine)
					}
				}
			} else {

				password := RandomString2(pwlength)
				hash, _ := HashPassword(password) // TODO add error handling

				fmt.Println(password, "	", hash)
				if testhash != false {
					matchLine = matchPasswordAndHash(password, hash)
					hashlines = append(hashlines, matchLine)
				}
			}

			count -= 1
		}
	} else {

		if infile != "" {
			log.Fatalf("Cannot run -s and -f together please select one")
		}
		password := thispw
		hash, _ := HashPassword(password)

		// fmt.Println(password, "	", hash)
		hasharray := []string{password,hash}
		hashline := strings.Join(hasharray, "	")
		fmt.Println(hashline)
		if testhash != false {
			match := checkHashAndPassword(password, hash)
			fmt.Println("Match:   ", match)
		}
	}

}
