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
	var fileErrors sfv.ErrorSummary

	scanner, err := sfv.NewFileScanner(filename)
	if err != nil {
		return err
	}

	for {
		for scanner.Scan() {
			entry := scanner.Entry()
			crc32, err := sfv.Crc32File(entry.Filename)

			if err != nil {
				fileErrors.Add(err)
				log.Print(err)
				continue
			}

			if entry.ExpectedCrc != crc32 {
				err := sfv.ErrMismatch{
					Filename:    entry.Filename,
					ExpectedCrc: entry.ExpectedCrc,
					ActualCrc:   crc32,
				}
				log.Print(err)
				fileErrors.Add(err)
				continue
			}

			if !quiet {
				log.Printf("%s: OK", entry.Filename)
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

	return fileErrors.Summary()
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
