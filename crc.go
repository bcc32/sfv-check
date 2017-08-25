package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

type errMismatch struct {
	filename    string
	expectedCrc uint32
	actualCrc   uint32
}

func (this errMismatch) Error() string {
	return fmt.Sprintf(
		"%s: NOT OK, expected %08X got %08X",
		this.filename,
		this.expectedCrc,
		this.actualCrc,
	)
}

type errFileOpen struct {
	filename string
}

func (this errFileOpen) Error() string {
	return fmt.Sprintf(
		"%s: NOT OK, file could not be read",
		this.filename,
	)
}

func checkFile(filename string, expectedCrc uint32) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return errFileOpen{filename}
	}

	hash := crc32.NewIEEE()

	io.Copy(hash, file)

	actualCrc := hash.Sum32()

	if expectedCrc != actualCrc {
		return errMismatch{filename, expectedCrc, actualCrc}
	}

	return nil
}
