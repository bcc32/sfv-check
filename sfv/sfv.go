/*

Package sfv contains code for parsing and checking SFV (Simple File
Verification) checksum files. SFV files are used to check for file corruption,
but do not prove file authenticity, i.e., the check is not cryptographically
secure.

https://en.wikipedia.org/wiki/Simple_file_verification

*/
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

// Entry represents an SFV line, consisting of the named file and its expected
// CRC-32 checksum.
type Entry struct {
	Filename    string
	ExpectedCRC uint32
}

// Check calculates the CRC-32 checksum of the named file and compares it to the
// expected checksum, returning a Result characterizing the outcome.
func (e Entry) Check() Result {
	actualCRC, err := CRC32File(e.Filename)
	if err != nil {
		return errResult{err, e.Filename}
	}
	if e.ExpectedCRC != actualCRC {
		return errMismatch{e.Filename, e.ExpectedCRC, actualCRC}
	}
	return okResult{e.Filename}
}

// TODO move result related code to another file, including err-summary.go

// A Result represents the result of checking a single SFV entry.
type Result interface {
	fmt.Stringer // format the Result like md5sum(1) and co.
	TAP() string // format the Result as a line of TAP (Test Anything Protocol)
	Err() error  // nil if file exists and matches checksum
}

// okResult represents a file that exists and matches its expected checksum.
type okResult struct {
	filename string
}

func (r okResult) String() string {
	return fmt.Sprintf("%s: OK", r.filename)
}

func (r okResult) TAP() string {
	return fmt.Sprintf("ok %s", r.filename)
}

func (r okResult) Err() error {
	return nil
}

// errResult represents an error that occurred during the calculation of the
// CRC-32 checksum of the named file.
type errResult struct {
	error
	filename string
}

func (r errResult) String() string {
	return fmt.Sprintf("%s: ERROR %s", r.filename, r.error)
}

func (r errResult) TAP() string {
	return fmt.Sprintf("not ok %s %s", r.filename, r.error)
}

func (r errResult) Err() error {
	return r
}

// errMismatch represents a mismatch between the expected and actual CRC-32
// checksums of the named file.
type errMismatch struct {
	Filename    string
	ExpectedCRC uint32
	ActualCRC   uint32
}

func (e errMismatch) String() string {
	return fmt.Sprintf(
		"%s: NOT OK, %s",
		e.Filename,
		e.Error(),
	)
}

func (e errMismatch) Error() string {
	return fmt.Sprintf(
		"expected %08X got %08X",
		e.ExpectedCRC,
		e.ActualCRC,
	)
}

func (e errMismatch) TAP() string {
	return fmt.Sprintf(
		"not ok expected %08X got %08X file %s",
		e.ExpectedCRC,
		e.ActualCRC,
		e.Filename,
	)
}

func (e errMismatch) Err() error {
	return e
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
	entry.ExpectedCRC = uint32(crc)

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
