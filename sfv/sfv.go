package sfv

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var Quiet = false

var ErrMalformedSfvLine = errors.New("malformed SFV line")

type ErrMismatch struct {
	filename    string
	expectedCrc uint32
	actualCrc   uint32
}

func (this ErrMismatch) Error() string {
	return fmt.Sprintf(
		"%s: NOT OK, expected %08X got %08X",
		this.filename,
		this.expectedCrc,
		this.actualCrc,
	)
}

type Entry struct {
	filename    string
	expectedCrc uint32
}

func parseSfvLine(line string) (entry Entry, err error) {
	if len(line) < 8 {
		err = ErrMalformedSfvLine
		return
	}

	filename, hex := line[:len(line)-8], line[len(line)-8:]
	entry.filename = strings.TrimSpace(filename)

	if entry.filename == "" {
		err = ErrMalformedSfvLine
		return
	}

	crc, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		err = ErrMalformedSfvLine
		return
	}
	entry.expectedCrc = uint32(crc)

	return
}

func CheckSfvFile(sfvFilename string) error {
	file, err := os.OpenFile(sfvFilename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	var fileErrors errorSummary

	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		if strings.HasPrefix(scanner.Text(), ";") {
			continue
		}
		entry, err := parseSfvLine(scanner.Text())
		if err != nil {
			log.Fatalf("%s:%d: %s", sfvFilename, lineNumber, err)
		}

		crc32, err := crc32File(entry.filename)
		if err != nil {
			fileErrors = append(fileErrors, err)
			log.Print(err)
		} else {
			if entry.expectedCrc != crc32 {
				log.Printf(
					"%s: NOT OK, expected %08X but got %08X",
					entry.filename,
					entry.expectedCrc,
					crc32,
				)
			} else if !Quiet {
				log.Printf("%s: OK", entry.filename)
			}
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
