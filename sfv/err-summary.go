package sfv

import (
	"errors"
	"fmt"
	"strings"
)

var errEmptyMultiError = errors.New("Error() called on empty multiError")

type errorSummary []error

func (this errorSummary) Error() string {
	mismatches := 0
	fileErrs := 0
	for _, e := range this {
		if e != nil {
			if _, ok := e.(ErrMismatch); ok {
				mismatches++
			}
			if _, ok := e.(ErrFileOpen); ok {
				fileErrs++
			}
		}
	}

	var lines []string
	if mismatches > 0 {
		msg := fmt.Sprintf("%d file(s) did NOT match", mismatches)
		lines = append(lines, msg)
	}
	if fileErrs > 0 {
		msg := fmt.Sprintf("%d file(s) could NOT be read", fileErrs)
		lines = append(lines, msg)
	}

	if len(lines) == 0 {
		panic(errEmptyMultiError)
	}

	return strings.Join(lines, "\n")
}
