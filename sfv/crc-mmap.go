// +build linux darwin

package sfv

import (
	"hash/crc32"
	"syscall"
)

func Crc32File(filename string) (uint32, error) {
	fd, err := syscall.Open(filename, syscall.O_RDONLY, 0)
	if err != nil {
		return 0, errFileOpen{err, filename}
	}

	var stat syscall.Stat_t
	err = syscall.Fstat(fd, &stat)
	if err != nil {
		return 0, errFileOpen{err, filename}
	}

	buf, err := syscall.Mmap(fd, 0, int(stat.Size),
		syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return 0, errFileOpen{err, filename}
	}

	hash := crc32.NewIEEE()
	hash.Write(buf)

	return hash.Sum32(), nil
}
