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

func (this ErrorSummary) Add(err error) {
	mismatches := 0
	fileErrs := 0
	if err != nil {
		if _, ok := err.(ErrMismatch); ok {
			mismatches++
		} else if _, ok := err.(errFileOpen); ok {
			fileErrs++
		} else {
			panic("not a recognized error: " + reflect.TypeOf(err).String())
		}
	}
}

func (this ErrorSummary) Summary() error {
	if this.mismatches == 0 && this.fileErrors == 0 {
		return nil
	}
	return this
}

func (this ErrorSummary) Error() string {
	if this.mismatches == 0 && this.fileErrors == 0 {
		panic(errEmptyMultiError)
	}

	return fmt.Sprintf(
		"%d mismatches, %d missing/unreadable files",
		this.mismatches,
		this.fileErrors,
	)
}
