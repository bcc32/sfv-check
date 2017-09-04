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

var errMalformedSFVLine = errors.New("malformed SFV line")

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
	fmt.Stringer      // format like md5sum(1) and co.
	TAP(i int) string // format as a line of TAP (Test Anything Protocol)
	Err() error       // nil if file exists and matches checksum
}

// okResult represents a file that exists and matches its expected checksum.
type okResult struct {
	filename string
}

func (r okResult) String() string {
	return fmt.Sprintf("%s: OK", r.filename)
}

func (r okResult) TAP(i int) string {
	return fmt.Sprintf("ok %d - %s", i, r.filename)
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

func (r errResult) TAP(i int) string {
	return fmt.Sprintf("not ok %d - %s %s", i, r.filename, r.error)
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

func (e errMismatch) TAP(i int) string {
	return fmt.Sprintf(
		"not ok %d - expected %08X got %08X file %s",
		i,
		e.ExpectedCRC,
		e.ActualCRC,
		e.Filename,
	)
}

func (e errMismatch) Err() error {
	return e
}

func parseSFVLine(line string) (entry Entry, err error) {
	if len(line) < 8 {
		err = errMalformedSFVLine
		return
	}

	filename, hex := line[:len(line)-8], line[len(line)-8:]
	entry.Filename = strings.TrimSpace(filename)

	if entry.Filename == "" {
		err = errMalformedSFVLine
		return
	}

	crc, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		err = errMalformedSFVLine
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
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !stat.Mode().IsRegular() {
		return nil, errors.New("not a regular file: " + filename)
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

		entry, err := parseSFVLine(line)
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

// ReadAll reads in the entire named SFV file at once, returning a slice of
// entries. Returns the first error instead if one is encountered.
func ReadAll(filename string) ([]Entry, error) {
	fs, err := NewFileScanner(filename)
	if err != nil {
		return nil, err
	}

	var entries []Entry

	for fs.Scan() {
		entries = append(entries, fs.Entry())
	}

	if err := fs.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
