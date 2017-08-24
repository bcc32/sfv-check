package main

import (
	"strings"
)

type multiError []error

func (this multiError) Error() string {
	var lines []string
	for _, e := range this {
		if e != nil {
			lines = append(lines, e.Error())
		}
	}
	return strings.Join(lines, "\n")
}
