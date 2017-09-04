package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/bcc32/sfv-check/sfv"
)

var errEmptyMultiError = errors.New("Error() called on empty multiError")

// errorSummary represents the aggregate errors encountered while checking the
// entries in an SFV file.
type errorSummary struct {
	mismatches int
	fileErrors int
}

// Add increments the appropriate errorSummary counters based on the type of the
// argument. If the argument is not a recognized error type, Add panics.
func (e *errorSummary) Add(err error) {
	if err != nil {
		if _, ok := err.(sfv.ErrMismatch); ok {
			e.mismatches++
		} else if _, ok := err.(sfv.ErrFileOpen); ok {
			e.fileErrors++
		} else {
			panic("not a recognized error: " + reflect.TypeOf(err).String())
		}
	}
}

func (e errorSummary) empty() bool {
	return e.mismatches == 0 && e.fileErrors == 0
}

// Summary returns an error value that is either nil if the errorSummary is
// empty (zero), or the errorSummary itself otherwise. This should be called
// prior to calling Error.
func (e errorSummary) Summary() error {
	if e.empty() {
		return nil
	}
	return e
}

func (e errorSummary) Error() string {
	if e.empty() {
		panic(errEmptyMultiError)
	}

	return fmt.Sprintf(
		"%d bad CRCs, %d not found",
		e.mismatches,
		e.fileErrors,
	)
}
