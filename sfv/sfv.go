package sfv

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var errMalformedSfvLine = errors.New("malformed SFV line")

type errParse struct {
	err         error
	sfvFilename string
	lineNumber  int
}

func (this errParse) Error() string {
	return fmt.Sprintf(
		"%s:%d: %s",
		this.sfvFilename,
		this.lineNumber,
		this.err,
	)
}

type ErrMismatch struct {
	Filename    string
	ExpectedCrc uint32
	ActualCrc   uint32
}

func (this ErrMismatch) Error() string {
	return fmt.Sprintf(
		"%s: NOT OK, expected %08X got %08X",
		this.Filename,
		this.ExpectedCrc,
		this.ActualCrc,
	)
}

type Entry struct {
	Filename    string
	ExpectedCrc uint32
}

func parseSfvLine(line string) (entry Entry, err error) {
	if len(line) < 8 {
		err = errMalformedSfvLine
		return
	}

	filename, hex := line[:len(line)-8], line[len(line)-8:]
	entry.Filename = strings.TrimSpace(filename)

	if entry.Filename == "" {
		err = errMalformedSfvLine
		return
	}

	crc, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		err = errMalformedSfvLine
		return
	}
	entry.ExpectedCrc = uint32(crc)

	return
}

type SfvFileScanner struct {
	input      *bufio.Scanner
	filename   string
	entry      Entry
	err        error
	lineNumber int
}

func NewSfvFileScanner(filename string) (*SfvFileScanner, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &SfvFileScanner{
		input:    bufio.NewScanner(file),
		filename: filename,
	}, nil
}

func (this *SfvFileScanner) Scan() bool {
	for this.input.Scan() {
		line := this.input.Text()
		this.lineNumber++

		if strings.HasPrefix(line, ";") {
			continue
		}

		entry, err := parseSfvLine(line)
		if err != nil {
			this.err = errParse{err, this.filename, this.lineNumber}
			return false
		}

		this.entry = entry
		return true
	}

	return false
}

func (this SfvFileScanner) Entry() Entry {
	return this.entry
}

func (this SfvFileScanner) Err() error {
	err := this.err
	this.err = nil
	return err
}
