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

func (this *ErrorSummary) Add(err error) {
	if err != nil {
		if _, ok := err.(ErrMismatch); ok {
			this.mismatches++
		} else if _, ok := err.(errFileOpen); ok {
			this.fileErrors++
		} else {
			panic("not a recognized error: " + reflect.TypeOf(err).String())
		}
	}
}

func (this *ErrorSummary) Summary() error {
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
		"%d bad CRCs, %d not found",
		this.mismatches,
		this.fileErrors,
	)
}
