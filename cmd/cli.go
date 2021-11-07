package main

import (
	"fmt"

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
			Name:     "number",
			Aliases:  []string{"n"},
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
			Aliases:  []string{"i"},
			Usage:    "Specify a file to read passwords from",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "output-file",
			Aliases:  []string{"o"},
			Usage:    "Specify a file to write out the pass/hash to",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "cost",
			Aliases:  []string{"c"},
			Usage:    "Specify the cost to use (Min: 4, Max: 31) (default: 14)",
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
	cost := c.Int("cost")

	if c.String("password") != "" {
		hashLines, err := generateHashForPassword(c.String("password"), cost)
		if err != nil {
			fmt.Errorf("There was an error creating a hash: %s\n", err)
		}

		if err := returnOutput(hashLines, c.String("out-file")); err != nil {
			fmt.Errorf("There was an error writing lines to file: %s\n", err)
		}
		return nil
	}

	if c.String("input-file") != "" {
		lines, err := readLines(c.String("input-file"))

		if err != nil {
			fmt.Errorf("There was an error reading lines from file: %s\n", err)
		}

		var hashLines []string
		for _, password := range lines {
			hashLine, err := generateHashForPassword(password, cost)
			if err != nil {
				fmt.Errorf("There was an error creating a hash: %s\n", err)
			}
			hashLines = append(hashLines, hashLine...)
		}

		if err := returnOutput(hashLines, c.String("out-file")); err != nil {
			fmt.Errorf("There was an error writing lines to file: %s\n", err)
		}
		return nil
	}

	var hashLines []string
	pwLength := c.Int("length")
	for number := c.Int("number"); number > 0; number-- {
		password := randomString(pwLength)
		hashLine, err := generateHashForPassword(password, cost)
		if err != nil {
			fmt.Errorf("There was an error creating a hash: %s\n", err)
		}
		hashLines = append(hashLines, hashLine...)
	}

	if err := returnOutput(hashLines, c.String("out-file")); err != nil {
		fmt.Errorf("There was an error writing lines to file: %s\n", err)
	}

	return nil
}
