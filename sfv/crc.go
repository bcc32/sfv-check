package sfv

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

type ErrFileOpen struct {
	filename string
}

func (this ErrFileOpen) Error() string {
	return fmt.Sprintf(
		"%s: NOT OK, file could not be read",
		this.filename,
	)
}

func crc32File(filename string) (uint32, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return 0, ErrFileOpen{filename}
	}

	hash := crc32.NewIEEE()

	io.Copy(hash, file)

	return hash.Sum32(), nil
}
