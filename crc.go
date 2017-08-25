package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strings"
)

type errMismatch struct {
	filename    string
	expectedCrc string
	actualCrc   string
}

func (this errMismatch) Error() string {
	return fmt.Sprintf(
		"%s: NOT OK, expected %s got %s",
		this.filename,
		strings.ToUpper(this.expectedCrc),
		strings.ToUpper(this.actualCrc),
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

func checkFile(filename string, expectedCrc string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return errFileOpen{filename}
	}

	hash := crc32.NewIEEE()

	io.Copy(hash, file)

	crc32 := hash.Sum32()
	actualCrc := fmt.Sprintf("%08X", crc32)

	if !strings.EqualFold(expectedCrc, actualCrc) {
		return errMismatch{filename, expectedCrc, actualCrc}
	}

	return nil
}
