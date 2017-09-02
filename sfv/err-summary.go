package sfv

import (
	"errors"
	"fmt"
	"reflect"
)

var errEmptyMultiError = errors.New("Error() called on empty multiError")

type ErrorSummary []error

func (this ErrorSummary) Error() string {
	mismatches := 0
	fileErrs := 0
	for _, e := range this {
		if e != nil {
			if _, ok := e.(ErrMismatch); ok {
				mismatches++
			} else if _, ok := e.(errFileOpen); ok {
				fileErrs++
			} else {
				panic("not a recognized error: " + reflect.TypeOf(e).String())
			}
		}
	}

	if mismatches == 0 && fileErrs == 0 {
		panic(errEmptyMultiError)
	}

	return fmt.Sprintf(
		"%d mismatches, %d missing/unreadable files",
		mismatches,
		fileErrs,
	)
}
