package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var quiet bool

func init() {
	log.SetFlags(0)

	const (
		defaultQuiet = false
		usageQuiet   = "suppress OK output for each correct file"
	)

	flag.BoolVar(&quiet, "quiet", defaultQuiet, usageQuiet)
	flag.BoolVar(&quiet, "q", defaultQuiet, usageQuiet+" (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			"Usage: %s [options] SFV-FILE [SFV-FILE]...\n",
			os.Args[0],
		)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	sfvFiles := flag.Args()

	success := true

	for _, file := range sfvFiles {
		err := checkSfvFile(file)

		if err != nil {
			success = false
			log.Printf("%s: %s\n", os.Args[0], err)
		}
	}

	if !success {
		os.Exit(1)
	}
}
