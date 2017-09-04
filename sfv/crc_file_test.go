package sfv

import (
	"io/ioutil"
	"os"
	"testing"
)

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func tempFileCrc32(data []byte) (uint32, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer func() {
		maybePanic(f.Close())
		maybePanic(os.Remove(f.Name()))
	}()

	f.Write(data)

	return Crc32File(f.Name())
}

func testCrc32(t *testing.T, expected uint32, data []byte) {
	t.Helper()

	actual, err := tempFileCrc32(data)
	if err != nil {
		t.Fatal(err)
	}
	if actual != expected {
		t.Fatalf("expected %08X, actual %08X", expected, actual)
	}
}

func TestCrc32File_empty(t *testing.T) {
	testCrc32(t, 0, []byte{})
}

func TestCrc32File_nonempty(t *testing.T) {
	testCrc32(t, 0x7C9CA35A, []byte{0xDE, 0xAD, 0xBE, 0xEF})
}
