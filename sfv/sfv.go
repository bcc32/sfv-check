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

func (e errParse) Error() string {
	return fmt.Sprintf(
		"%s:%d: %s",
		e.sfvFilename,
		e.lineNumber,
		e.err,
	)
}

type ErrMismatch struct {
	Filename    string
	ExpectedCrc uint32
	ActualCrc   uint32
}

func (e ErrMismatch) Error() string {
	return fmt.Sprintf(
		"%s: NOT OK, expected %08X got %08X",
		e.Filename,
		e.ExpectedCrc,
		e.ActualCrc,
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

type FileScanner struct {
	input      *bufio.Scanner
	filename   string
	entry      Entry
	err        error
	lineNumber int
}

func NewSfvFileScanner(filename string) (*FileScanner, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &FileScanner{
		input:    bufio.NewScanner(file),
		filename: filename,
	}, nil
}

func (fs *FileScanner) Scan() bool {
	for fs.input.Scan() {
		line := fs.input.Text()
		fs.lineNumber++

		if strings.HasPrefix(line, ";") {
			continue
		}

		entry, err := parseSfvLine(line)
		if err != nil {
			fs.err = errParse{err, fs.filename, fs.lineNumber}
			return false
		}

		fs.entry = entry
		return true
	}

	return false
}

func (fs *FileScanner) Entry() Entry {
	return fs.entry
}

func (fs *FileScanner) Err() error {
	err := fs.err
	fs.err = nil
	return err
}
