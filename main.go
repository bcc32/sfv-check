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
var tap bool

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	const (
		defaultQuiet = false
		usageQuiet   = "suppress OK output for each correct file"

		defaultTap = false
		usageTap   = "print results in TAP format"
	)

	flag.BoolVar(&quiet, "quiet", defaultQuiet, usageQuiet)
	flag.BoolVar(&quiet, "q", defaultQuiet, usageQuiet+" (shorthand)")

	flag.BoolVar(&tap, "tap", defaultTap, usageTap)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			"Usage: %s [options] SFV-FILE [SFV-FILE]...\n",
			os.Args[0],
		)
		flag.PrintDefaults()
	}
}

func displayResult(r sfv.Result) {
	if tap {
		log.Print(r.TAP())
		return
	}
	if !quiet || r.Err() != nil {
		// FIXME this is a bit precarious, since missing the call to
		// String() would result in Error() being called.
		log.Print(r.String())
	}
}

func checkSFVFile(filename string, results *sfv.ErrorSummary) error {
	scanner, err := sfv.NewFileScanner(filename)
	if err != nil {
		return err
	}

	for {
		for scanner.Scan() {
			entry := scanner.Entry()

			result := entry.Check()
			results.Add(result)
			displayResult(result)
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

	var results sfv.ErrorSummary

	for _, file := range sfvFiles {
		err := checkSFVFile(file, &results)
		if err != nil {
			log.Print(err)
		}
	}

	if err := results.Summary(); err != nil {
		log.Printf("%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}
