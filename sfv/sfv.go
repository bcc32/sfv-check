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

// ErrMismatch represents a mismatch between the expected and actual CRC-32
// checksums of the named file.
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

// Entry represents an SFV line, consisting of the named file and its expected
// CRC-32 checksum.
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

// A FileScanner parses an SFV file, reporting any syntax or I/O errors
// encountered.
type FileScanner struct {
	input      *bufio.Scanner
	filename   string
	entry      Entry
	err        error
	lineNumber int
}

// NewFileScanner constructs a new FileScanner, returning an error if the named
// file cannot be opened.
func NewFileScanner(filename string) (*FileScanner, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &FileScanner{
		input:    bufio.NewScanner(file),
		filename: filename,
	}, nil
}

// Scan scans through the file for the next SFV entry, returning true if it
// finds one. It returns false if an I/O or syntax error is encountered.
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

	fs.err = nil
	return false
}

// Entry returns the last entry parsed by Scan. It should only be called after a
// call to Scan returns true.
func (fs *FileScanner) Entry() Entry {
	return fs.entry
}

// Err returns the last error encountered by Scan. It should only be called
// after a call to Scan returns false. If Err returns nil, then the end of the
// file has been reached without error.
func (fs *FileScanner) Err() error {
	return fs.err
}
