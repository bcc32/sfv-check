package main

import (
	"errors"
	"fmt"
	"strings"
)

var errEmptyMultiError = errors.New("Error() called on empty multiError")

type multiError []error

func (this multiError) Error() string {
	mismatches := 0
	fileErrs := 0
	for _, e := range this {
		if e != nil {
			if _, ok := e.(errMismatch); ok {
				mismatches++
			}
			if _, ok := e.(errFileOpen); ok {
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
