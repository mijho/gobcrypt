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
	Usage:       "Password to hash",
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

var costValue int
var costFlag = &cli.IntFlag{
	Name:        "cost",
	Aliases:     []string{"c"},
	Usage:       "Specify the cost to use (Min: 4, Max: 31) (default: 14)",
	Required:    false,
	Destination: &costValue,
}

var CLIApp = &cli.App{
	Name:        "gobcrypt",
	Usage:       "A Bcrypt hash/password generator",
	ArgsUsage:   "",
	Version:     Version + "-" + Commit,
	Description: fmt.Sprintf("Build: %s", Build),
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

// ValidateHandler provides functionality to validate pass/hash pairs
func ValidateHandler(c *cli.Context) error {
	var hashLines []string

	if inputFileValue != "" {
		lines, err := readLines(inputFileValue)

		if err != nil {
			fmt.Errorf("There was an error reading lines from file: %s\n", err)
		}
		for _, line := range lines {
			result, _ := matchPasswordAndHash(line)
			hashLines = append(hashLines, result)
		}
		if err := returnOutput(hashLines, outputFileValue); err != nil {
			fmt.Errorf("There was an error writing lines to file: %s\n", err)
		}		
	} else {
		hashLine := strings.Join([]string{passwordValue, hashValue}, " ")
		result, _ := matchPasswordAndHash(hashLine)
		hashLines = append(hashLines, result)

		if err := returnOutput(hashLines, outputFileValue); err != nil {
			fmt.Errorf("There was an error writing lines to file: %s\n", err)
		}
	}

	return nil
}

// GenerateHandler provides functionality to generate pass/hash pairs
func GenerateHandler(c *cli.Context) error {
	cost := costValue

	if passwordValue != "" {
		hashLines, err := generateHashForPassword(passwordValue, cost)
		if err != nil {
			fmt.Errorf("There was an error creating a hash: %s\n", err)
		}

		if err := returnOutput(hashLines, outputFileValue); err != nil {
			fmt.Errorf("There was an error writing lines to file: %s\n", err)
		}
		return nil
	}

	if inputFileValue != "" {
		lines, err := readLines(inputFileValue)

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

		if err := returnOutput(hashLines, outputFileValue); err != nil {
			fmt.Errorf("There was an error writing lines to file: %s\n", err)
		}
		return nil
	}

	var hashLines []string
	pwLength := lengthValue
	for number := numberValue; number > 0; number-- {
		password := randomString(pwLength)
		hashLine, err := generateHashForPassword(password, cost)
		if err != nil {
			fmt.Errorf("There was an error creating a hash: %s\n", err)
		}
		hashLines = append(hashLines, hashLine...)
	}

	if err := returnOutput(hashLines, outputFileValue); err != nil {
		fmt.Errorf("There was an error writing lines to file: %s\n", err)
	}

	return nil
}
