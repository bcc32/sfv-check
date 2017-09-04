// Command sfv-check accepts one or more SFV-formatted checksum files, and
// verifies the contents of the files listed therein. Mismatches, file errors,
// and success notifications are printed to standard output.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bcc32/sfv-check/sfv"
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

func checkSfvFile(filename string) error {
	var results sfv.ErrorSummary

	scanner, err := sfv.NewFileScanner(filename)
	if err != nil {
		return err
	}

	for {
		for scanner.Scan() {
			entry := scanner.Entry()
			result := entry.Check()

			results.Add(result)

			if !quiet || result.Err() != nil {
				// FIXME this is a bit precarious, since missing the call to
				// String() would result in Error() being called.
				log.Print(result.String())
			}
		}
		if err := scanner.Err(); err != nil {
			log.Print(err)
		} else {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return results.Summary()
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
