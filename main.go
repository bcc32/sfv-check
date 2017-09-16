// Command sfv-check accepts one or more SFV-formatted checksum files, and
// verifies the contents of the files listed therein. Mismatches, file errors,
// and success notifications are printed to standard output.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/bcc32/sfv-check/sfv"
)

var quiet bool
var tap bool
var changeDir bool

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	const (
		defaultQuiet = false
		usageQuiet   = "suppress OK output for each correct file"

		defaultTap = false
		usageTap   = "print results in TAP format (one SFV file only)"

		defaultChangeDir = false
		usageChangeDir   = "chdir to the directory containing each SFV file"
	)

	flag.BoolVar(&quiet, "quiet", defaultQuiet, usageQuiet)
	flag.BoolVar(&quiet, "q", defaultQuiet, usageQuiet+" (shorthand)")

	flag.BoolVar(&tap, "tap", defaultTap, usageTap)

	flag.BoolVar(&changeDir, "cd", defaultChangeDir, usageChangeDir)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			"Usage: %s [options] SFV-FILE [SFV-FILE]...\n",
			os.Args[0],
		)
		flag.PrintDefaults()
	}
}

// checkSFVFile returns an error only if the SFV file cannot be read or is
// malformed. File mismatches and errors are recorded in results.
func checkSFVFile(filename string, results *sfv.ResultSummary) error {
	scanner, err := sfv.NewFileScanner(filename)
	if err != nil {
		return err
	}

	if changeDir {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		os.Chdir(path.Dir(filename))
		defer os.Chdir(dir)
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

	return nil
}

// tapSFVFile returns an error only if the SFV file cannot be read or is
// malformed. File mismatches and errors are recorded in results.
func tapSFVFile(filename string, results *sfv.ResultSummary) error {
	entries, err := sfv.ReadAll(filename)
	if err != nil {
		return err
	}

	if changeDir {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		os.Chdir(path.Dir(filename))
		defer os.Chdir(dir)
	}

	log.Printf("1..%d\n", len(entries))

	for i, entry := range entries {
		result := entry.Check()
		results.Add(result)
		log.Print(result.TAP(i + 1))
	}

	return nil
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	sfvFiles := flag.Args()

	var exitCode int
	var results sfv.ResultSummary

	if tap {
		if len(sfvFiles) != 1 {
			log.Fatal("-tap can only be used with one SFV file")
		}
		err := tapSFVFile(sfvFiles[0], &results)
		if err != nil {
			log.Print(err)
			exitCode = 1
		}
	} else {
		for _, file := range sfvFiles {
			err := checkSFVFile(file, &results)
			if err != nil {
				log.Print(err)
				exitCode = 1
			}
		}
	}

	if err := results.Summary(); err != nil {
		log.Printf("%s: %s\n", os.Args[0], err)
		exitCode = 1
	}

	os.Exit(exitCode)
}
