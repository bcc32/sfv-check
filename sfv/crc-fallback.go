// +build !linux,!darwin

package sfv

import (
	"hash/crc32"
	"io"
	"os"
)

func Crc32File(filename string) (crc uint32, error error) {
	file, err := os.Open(filename)
	if err != nil {
		error = err
		return
	}

	defer func() {
		err := file.Close()
		if err != nil {
			error = err
		}
	}()

	hash := crc32.NewIEEE()

	io.Copy(hash, file)
	crc = hash.Sum32()
	return
}
