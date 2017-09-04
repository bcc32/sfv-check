package sfv

import (
	"fmt"
)

type errFileOpen struct {
	filename string
}

func (this errFileOpen) Error() string {
	return fmt.Sprintf(
		"%s: NOT OK, file could not be read",
		this.filename,
	)
}
