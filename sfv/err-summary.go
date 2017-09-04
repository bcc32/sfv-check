package sfv

import (
	"errors"
	"fmt"
	"reflect"
)

var errEmptyMultiError = errors.New("Error() called on empty multiError")

type ErrorSummary struct {
	mismatches int
	fileErrors int
}

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

func (e *ErrorSummary) Summary() error {
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
