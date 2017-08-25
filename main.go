package main

import (
	"fmt"
	"log"
	"os"
)

var verbose bool = true // TODO

func init() {
	log.SetFlags(0)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s SFV-FILE", os.Args[0])
		os.Exit(1)
	}

	sfvFile := os.Args[1]
	err := checkSfvFile(sfvFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}
