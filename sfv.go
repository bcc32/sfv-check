package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

var errMalformedSfvLine = errors.New("malformed SFV line")

func parseSfvLine(line string) (filename string, expectedCrc uint32, err error) {
	if len(line) < 8 {
		err = errMalformedSfvLine
		return
	}

	filename, hex := line[:len(line)-8], line[len(line)-8:]
	filename = strings.TrimSpace(filename)

	if filename == "" {
		err = errMalformedSfvLine
		return
	}

	crc, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		err = errMalformedSfvLine
		return
	}
	expectedCrc = uint32(crc)

	return
}

func checkSfvFile(sfvFilename string) error {
	file, err := os.OpenFile(sfvFilename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	var fileErrors multiError

	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		if strings.HasPrefix(scanner.Text(), ";") {
			continue
		}
		filename, expectedCrc, err := parseSfvLine(scanner.Text())
		if err != nil {
			log.Fatalf("%s:%d: %s", sfvFilename, lineNumber, err)
		}

		err = checkFile(filename, expectedCrc)
		if err != nil {
			fileErrors = append(fileErrors, err)
			log.Print(err)
		} else if !quiet {
			log.Printf("%s: OK", filename)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// necessary because `fileErrors` is declared with a concrete type
	if fileErrors == nil {
		return nil
	}

	return fileErrors
}
