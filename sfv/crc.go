package sfv

import (
	"fmt"
)

// ErrFileOpen represents an error that occurred trying to open or read the
// named file.
type ErrFileOpen struct {
	err      error
	filename string
}

func (e ErrFileOpen) Error() string {
	return fmt.Sprintf(
		"%s: NOT OK, %s",
		e.filename,
		e.err,
	)
}
