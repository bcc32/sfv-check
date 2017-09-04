package sfv

import (
	"errors"
	"fmt"
)

var errEmptyMultiError = errors.New("Error() called on empty multiError")

// ResultSummary represents the aggregate errors encountered while checking the
// entries in an SFV file.
type ResultSummary struct {
	mismatches int
	fileErrors int
	totalTests int
}

// Add increments the appropriate errorSummary counters based on the type of the
// argument.
func (e *ResultSummary) Add(r Result) {
	e.totalTests++

	if err := r.Err(); err != nil {
		if _, ok := err.(errMismatch); ok {
			e.mismatches++
		} else {
			e.fileErrors++
		}
	}
}

func (e ResultSummary) empty() bool {
	return e.mismatches == 0 && e.fileErrors == 0
}

// Summary returns an error value that is either nil if the ResultSummary is
// empty (zero), or the ResultSummary itself otherwise. This should be called
// prior to calling Error.
func (e ResultSummary) Summary() error {
	if e.empty() {
		return nil
	}
	return e
}

func (e ResultSummary) Error() string {
	if e.empty() {
		panic(errEmptyMultiError)
	}

	return fmt.Sprintf(
		"%d bad CRCs, %d file errors",
		e.mismatches,
		e.fileErrors,
	)
}
