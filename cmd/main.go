package main

import (
	"log"
	"os"
)

func main() {
	err := CLIApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
