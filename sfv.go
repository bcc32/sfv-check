package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var malformedSfvLine = errors.New("malformed SFV line")

func parseSfvLine(line string) (filename, expectedCrc string, err error) {
	if len(line) < 8 {
		err = malformedSfvLine
		return
	}

	filename, expectedCrc = line[:len(line)-8], line[len(line)-8:]
	filename = strings.TrimSpace(filename)

	if filename == "" {
		err = malformedSfvLine
		return
	}

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
		filename, expectedCrc, err := parseSfvLine(scanner.Text())
		if err != nil {
			// TODO use log package
			fmt.Printf("%s:%d: %s\n", sfvFilename, lineNumber, err)
		}
		err = checkFile(filename, expectedCrc)
		if err != nil {
			fileErrors = append(fileErrors, err)
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
