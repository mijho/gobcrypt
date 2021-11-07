package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

// Build time variables
var (
	Build   string
	Commit  string
	Name    string
	Version string
	Author  string
	Email   string
)

var CLIApp = &cli.App{
	Name:        "gobcrypt",
	Usage:       "A Bcrypt hash/password generator",
	ArgsUsage:   "",
	Version:     Version + "-" + Commit,
	Description: fmt.Sprintf("Build: %s", Build),
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "count",
			Aliases:  []string{"c"},
			Usage:    "Specify the number of hashes to create (default: 1)",
			Value:    1,
			Required: false,
		},
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
			Usage:    "Hash the specified password only",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "length",
			Aliases:  []string{"l"},
			Usage:    "Specify the length of password required (default: 16)",
			Value:    16,
			Required: false,
		},
		&cli.StringFlag{
			Name:     "input-file",
			Aliases:  []string{"f"},
			Usage:    "Specify a file to read passwords from",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "output-file",
			Aliases:  []string{"o"},
			Usage:    "Specify a file to write out the pass/hash to",
			Required: false,
		},
	},
	Action: LocalHandler,
	Authors: []*cli.Author{
		{
			Name:  Author,
			Email: Email,
		},
	},
}

// LocalHandler provides stand alone functionality to generate CDX from the
// provided WARC locally
func LocalHandler(c *cli.Context) error {

	if c.String("password") != "" {
		hashLines, err := generateHashForPassword(c.String("password"))
		if err != nil {
			fmt.Errorf("There was an error creating a hash: %s\n", err)
		}

		if c.String("output-file") != "" {
			if err := writeLines(hashLines, c.String("output-file")); err != nil {
				fmt.Errorf("There was an error writing lines to file: %s\n", err)
			} else {
				fmt.Println(hashLines[0])
			}
		}
		return nil
	}

	



	return nil
}

func generateHashForPassword(password string) ([]string, error) {
	fmt.Printf("Generating hash for password: %s\n", password)
	var hashLines []string

	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	hashLine := strings.Join([]string{password, hash}, " ")
	hashLines = append(hashLines, hashLine)

	return hashLines, err
}
