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

var numberValue int
var numberFlag = &cli.IntFlag{
	Name:        "number",
	Aliases:     []string{"n"},
	Usage:       "Specify the number of hashes to create (default: 1)",
	Value:       1,
	Required:    false,
	Destination: &numberValue,
}

var passwordValue string
var passwordFlag = &cli.StringFlag{
	Name:        "password",
	Aliases:     []string{"p"},
	Usage:       "Password to hash (generate) or verify (validate)",
	Required:    false,
	Destination: &passwordValue,
}

var hashValue string
var hashFlag = &cli.StringFlag{
	Name:        "hash",
	Aliases:     []string{"hp"},
	Usage:       "Hash to verify",
	Required:    false,
	Destination: &hashValue,
}

var lengthValue int
var lengthFlag = &cli.IntFlag{
	Name:        "length",
	Aliases:     []string{"l"},
	Usage:       "Specify the length of password required (default: 16)",
	Value:       16,
	Required:    false,
	Destination: &lengthValue,
}

var inputFileValue string
var inputFileFlag = &cli.StringFlag{
	Name:        "input-file",
	Aliases:     []string{"i"},
	Usage:       "Specify a file to read passwords from",
	Required:    false,
	Destination: &inputFileValue,
}

var outputFileValue string
var outputFileFlag = &cli.StringFlag{
	Name:        "output-file",
	Aliases:     []string{"o"},
	Usage:       "Specify a file to write out the pass/hash to",
	Required:    false,
	Destination: &outputFileValue,
}

// defaultCost matches the default documented in the README.
const defaultCost = 14

var costValue int
var costFlag = &cli.IntFlag{
	Name:        "cost",
	Aliases:     []string{"c"},
	Usage:       "Specify the cost to use (Min: 4, Max: 31) (default: 14)",
	Value:       defaultCost,
	Required:    false,
	Destination: &costValue,
}

// CLIApp constructs the cli applications
var CLIApp = &cli.App{
	Name:        "gobcrypt",
	Usage:       "A Bcrypt hash/password generator",
	ArgsUsage:   "",
	Version:     Version + "-" + Commit,
	Description: fmt.Sprintf("%s: Build: %s", Name, Build),
	Commands: []*cli.Command{
		{
			Name:  "generate",
			Usage: "generate pass/hash pairs",
			Flags: []cli.Flag{
				numberFlag,
				passwordFlag,
				lengthFlag,
				inputFileFlag,
				outputFileFlag,
				costFlag,
			},
			Action: GenerateHandler,
		},
		{
			Name:  "validate",
			Usage: "validate pass/hash pairs",
			Flags: []cli.Flag{
				passwordFlag,
				hashFlag,
				inputFileFlag,
				outputFileFlag,
			},
			Action: ValidateHandler,
		},
	},
	Authors: []*cli.Author{
		{
			Name:  Author,
			Email: Email,
		},
	},
}

// ValidateHandler verifies pass/hash pairs. It exits non-zero if any pair fails
// to match or is malformed, so it can be used in scripts.
func ValidateHandler(c *cli.Context) error {
	var lines []string
	if inputFileValue != "" {
		var err error
		lines, err = readLines(inputFileValue)
		if err != nil {
			return cli.Exit(fmt.Sprintf("error reading lines from file: %s", err), 1)
		}
	} else {
		if passwordValue == "" || hashValue == "" {
			return cli.Exit("validate needs --password and --hash, or --input-file", 1)
		}
		lines = []string{passwordValue + " " + hashValue}
	}

	failed := 0
	results := make([]string, 0, len(lines))
	for _, line := range lines {
		result, err := matchPasswordAndHash(line)
		if err != nil {
			failed++
			if result == "" {
				result = fmt.Sprintf("MATCH: ERROR, line: %q, error: %s", line, err)
			}
		}
		results = append(results, result)
	}

	if err := returnOutput(results, outputFileValue); err != nil {
		return cli.Exit(fmt.Sprintf("error writing output: %s", err), 1)
	}
	if failed > 0 {
		return cli.Exit(fmt.Sprintf("%d of %d pair(s) did not validate", failed, len(lines)), 1)
	}
	return nil
}

// GenerateHandler generates pass/hash pairs from a given password, an input
// file, or random passwords.
func GenerateHandler(c *cli.Context) error {
	cost := costValue
	var hashLines []string

	switch {
	case passwordValue != "":
		lines, err := generateHashForPassword(passwordValue, cost)
		if err != nil {
			return cli.Exit(fmt.Sprintf("error creating a hash: %s", err), 1)
		}
		hashLines = lines

	case inputFileValue != "":
		passwords, err := readLines(inputFileValue)
		if err != nil {
			return cli.Exit(fmt.Sprintf("error reading lines from file: %s", err), 1)
		}
		for _, password := range passwords {
			lines, err := generateHashForPassword(password, cost)
			if err != nil {
				return cli.Exit(fmt.Sprintf("error creating a hash: %s", err), 1)
			}
			hashLines = append(hashLines, lines...)
		}

	default:
		for i := 0; i < numberValue; i++ {
			password, err := randomString(lengthValue)
			if err != nil {
				return cli.Exit(fmt.Sprintf("error generating random password: %s", err), 1)
			}
			lines, err := generateHashForPassword(password, cost)
			if err != nil {
				return cli.Exit(fmt.Sprintf("error creating a hash: %s", err), 1)
			}
			hashLines = append(hashLines, lines...)
		}
	}

	if err := returnOutput(hashLines, outputFileValue); err != nil {
		return cli.Exit(fmt.Sprintf("error writing output: %s", err), 1)
	}
	return nil
}
