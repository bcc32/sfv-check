package sfv

import (
	"fmt"
)

type errFileOpen struct {
	err      error
	filename string
}

func (this errFileOpen) Error() string {
	return fmt.Sprintf(
		"%s: NOT OK, %s",
		this.filename,
		this.err,
	)
}
