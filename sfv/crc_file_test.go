package sfv_test

import (
	"hash/crc32"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/bcc32/sfv-check/sfv"
)

const (
	megabyte = 1024 * 1024
)

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func tempFileCRC32(data []byte) (uint32, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer func() {
		maybePanic(f.Close())
		maybePanic(os.Remove(f.Name()))
	}()

	f.Write(data)

	return sfv.CRC32File(f.Name())
}

func testCRC32(t *testing.T, expected uint32, data []byte) {
	t.Helper()

	actual, err := tempFileCRC32(data)
	if err != nil {
		t.Error(err)
	}
	if actual != expected {
		t.Errorf("CRC-32(%X) = %08X; want %08X", data, actual, expected)
	}
}

func TestCRC32File_empty(t *testing.T) {
	testCRC32(t, 0, []byte{})
}

func TestCRC32File_nonempty(t *testing.T) {
	testCRC32(t, 0x7C9CA35A, []byte{0xDE, 0xAD, 0xBE, 0xEF})
}

func testCRC32NoFile(t *testing.T, filename string) {
	t.Helper()
	_, err := sfv.CRC32File(filename)
	if err == nil {
		t.Errorf("filename %q; expected error", filename)
	}
}

func TestCRC32File_noFile(t *testing.T) {
	testCRC32NoFile(t, "/zzzzzzzz")
}

func TestCRC32File_noFile_emptyName(t *testing.T) {
	testCRC32NoFile(t, "")
}

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func makeRandomFile(size int64) (filename string, crc uint32, error error) {
	file, err := ioutil.TempFile("", "")
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
	writer := io.MultiWriter(file, hash)

	_, err = io.CopyN(writer, random, size)
	if err != nil {
		error = err
		return
	}

	filename = file.Name()
	crc = hash.Sum32()
	return
}

func BenchmarkCRC32File_64MB(b *testing.B) {
	var bytes int64 = 64 * megabyte
	b.SetBytes(int64(bytes))

	filename, expected, err := makeRandomFile(bytes)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		maybePanic(os.Remove(filename))
	}()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		actual, err := sfv.CRC32File(filename)
		if err != nil {
			b.Fatal(err)
		}
		if actual != expected {
			b.Fatal("incorrect CRC-32")
		}
	}
}
