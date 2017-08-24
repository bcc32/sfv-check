package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strings"
)

func errorMismatch(filename, expectedCrc, actualCrc string) error {
	return fmt.Errorf(
		"%s: NOT OK, expected %s got %s",
		filename,
		expectedCrc,
		actualCrc,
	)
}

func checkFile(filename string, expectedCrc string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	hash := crc32.NewIEEE()

	io.Copy(hash, file)

	crc32 := hash.Sum32()
	actualCrc := fmt.Sprintf("%08X", crc32)

	if !strings.EqualFold(expectedCrc, actualCrc) {
		return errorMismatch(filename, expectedCrc, actualCrc)
	}

	return nil
}
