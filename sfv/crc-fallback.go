// +build !linux,!darwin

package sfv

import (
	"hash/crc32"
	"io"
	"os"
)

func Crc32File(filename string) (uint32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, errFileOpen{filename}
	}

	hash := crc32.NewIEEE()

	io.Copy(hash, file)

	return hash.Sum32(), nil
}
