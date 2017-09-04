package sfv

import (
	"hash/crc32"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

const (
	MEGABYTE = 1024 * 1024
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

func TestCrc32File_noFile(t *testing.T) {
	_, err := Crc32File("/zzzzzzzz")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCrc32File_noFile_emptyName(t *testing.T) {
	_, err := Crc32File("")
	if err == nil {
		t.Fatal("expected error")
	}
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

func BenchmarkCrc32File_64MB(b *testing.B) {
	var bytes int64 = 64 * MEGABYTE
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
		actual, err := Crc32File(filename)
		if err != nil {
			b.Fatal(err)
		}
		if expected != actual {
			b.Fatalf("expected %08X, got %08X", expected, actual)
		}
	}
}
