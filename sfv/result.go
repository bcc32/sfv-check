package sfv

import (
	"errors"
	"fmt"
)

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

// mismatchResult represents a mismatch between the expected and actual CRC-32
// checksums of the named file.
type mismatchResult struct {
	Filename    string
	ExpectedCRC uint32
	ActualCRC   uint32
}

func (r mismatchResult) String() string {
	return fmt.Sprintf(
		"%s: NOT OK, %s",
		r.Filename,
		r.Error(),
	)
}

func (r mismatchResult) Error() string {
	return fmt.Sprintf(
		"expected %08X got %08X",
		r.ExpectedCRC,
		r.ActualCRC,
	)
}

func (r mismatchResult) TAP(i int) string {
	return fmt.Sprintf(
		"not ok %d - expected %08X got %08X file %s",
		i,
		r.ExpectedCRC,
		r.ActualCRC,
		r.Filename,
	)
}

func (r mismatchResult) Err() error {
	return r
}

var errEmptyResultSummary = errors.New("Error() called on empty resultSummary")

// ResultSummary represents the aggregate errors encountered while checking the
// entries in an SFV file.
type ResultSummary struct {
	mismatches int
	fileErrors int
}

// Add increments the appropriate errorSummary counters based on the type of the
// argument.
func (rs *ResultSummary) Add(r Result) {
	if err := r.Err(); err != nil {
		if _, ok := err.(mismatchResult); ok {
			rs.mismatches++
		} else {
			rs.fileErrors++
		}
	}
}

func (rs ResultSummary) empty() bool {
	return rs.mismatches == 0 && rs.fileErrors == 0
}

// Summary returns an error value that is either nil if the ResultSummary is
// empty (zero), or the ResultSummary itself otherwise. This should be called
// prior to calling Error.
func (rs ResultSummary) Summary() error {
	if rs.empty() {
		return nil
	}
	return rs
}

func (rs ResultSummary) Error() string {
	if rs.empty() {
		panic(errEmptyResultSummary)
	}

	return fmt.Sprintf(
		"%d bad CRCs, %d file errors",
		rs.mismatches,
		rs.fileErrors,
	)
}
