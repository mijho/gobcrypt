// gobcrypt.go
package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
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
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	//password := "secret"
	//password := os.Args[1]
	var count, pwlength int
	var thispw string
	var testhash bool
	flag.IntVar(&count, "c", 1, "Specify the number of hashes to create")
	flag.StringVar(&thispw, "s", "", "Hash the specified password only (1)")
	flag.BoolVar(&testhash, "t", false, "Validate the hash & pass")
	flag.IntVar(&pwlength, "l", 15, "Specify the length of password required")
	flag.Parse()

	if thispw == "" {
		for count > 0 {

			password := RandomString2(pwlength)
			hash, _ := HashPassword(password) // ignore error for the sake of simplicity

			fmt.Println(password, "	", hash)

			if testhash != false {
				match := CheckPasswordHash(password, hash)
				fmt.Println("Match:   ", match)
			}
			count -= 1
		}
	} else {
		password := thispw
		hash, _ := HashPassword(password)

		fmt.Println(password, "	", hash)
		if testhash != false {
			match := CheckPasswordHash(password, hash)
			fmt.Println("Match:   ", match)
		}
	}

}
