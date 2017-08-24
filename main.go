package main

import (
	"fmt"
	"os"
)

func main() {
	sfvFile := os.Args[1]
	err := checkSfvFile(sfvFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}
