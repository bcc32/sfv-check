package sfv

import (
	"errors"
	"fmt"
	"reflect"
)

var errEmptyMultiError = errors.New("Error() called on empty multiError")

// ErrorSummary represents the aggregate errors encountered while checking the
// entries in an SFV file.
type ErrorSummary struct {
	mismatches int
	fileErrors int
}

// Add increments the appropriate ErrorSummary counters based on the type of the
// argument. If the argument is not a recognized error type, Add panics.
func (e *ErrorSummary) Add(err error) {
	if err != nil {
		if _, ok := err.(ErrMismatch); ok {
			e.mismatches++
		} else if _, ok := err.(errFileOpen); ok {
			e.fileErrors++
		} else {
			panic("not a recognized error: " + reflect.TypeOf(err).String())
		}
	}
}

func (e ErrorSummary) empty() bool {
	return e.mismatches == 0 && e.fileErrors == 0
}

// Summary returns an error value that is either nil if the ErrorSummary is
// empty (zero), or the ErrorSummary itself otherwise. This should be called
// prior to calling Error.
func (e ErrorSummary) Summary() error {
	if e.empty() {
		return nil
	}
	return e
}

func (e ErrorSummary) Error() string {
	if e.empty() {
		panic(errEmptyMultiError)
	}

	return fmt.Sprintf(
		"%d bad CRCs, %d not found",
		e.mismatches,
		e.fileErrors,
	)
}
