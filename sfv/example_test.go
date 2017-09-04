package sfv_test

import (
	"fmt"

	"github.com/bcc32/sfv-check/sfv"
)

func ExampleCRC32File() {
	crc, err := sfv.CRC32File("testdata/test-file.go")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%08X\n", crc)
	fmt.Printf("%d\n", crc)
	// Output:
	// A80579F6
	// 2818931190
}

func ExampleFileScanner() {
	fs, err := sfv.NewFileScanner("testdata/test.sfv")
	if err != nil {
		panic(err)
	}

	fmt.Println(fs.Scan())
	fmt.Printf("%+v\n", fs.Entry())
	// Output:
	// true
	// {Filename:test-file.go ExpectedCRC:2818931190}
}
