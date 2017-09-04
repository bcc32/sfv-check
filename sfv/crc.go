package sfv

import (
	"fmt"
)

type errFileOpen struct {
	err      error
	filename string
}

func (e errFileOpen) Error() string {
	return fmt.Sprintf(
		"%s: NOT OK, %s",
		e.filename,
		e.err,
	)
}
