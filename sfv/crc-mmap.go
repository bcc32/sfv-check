// +build linux darwin

package sfv

import (
	"hash/crc32"
	"syscall"
)

func Crc32File(filename string) (crc uint32, error error) {
	fd, err := syscall.Open(filename, syscall.O_RDONLY, 0)
	if err != nil {
		error = err
		return
	}
	defer func() {
		err := syscall.Close(fd)
		if err != nil {
			error = err
		}
	}()

	var stat syscall.Stat_t
	err = syscall.Fstat(fd, &stat)
	if err != nil {
		error = err
		return
	}

	if stat.Size == 0 {
		return
	}

	buf, err := syscall.Mmap(fd, 0, int(stat.Size),
		syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		error = err
		return
	}
	defer func() {
		err := syscall.Munmap(buf)
		if err != nil {
			error = err
		}
	}()

	hash := crc32.NewIEEE()
	_, err = hash.Write(buf)
	if err != nil {
		error = err
		return
	}

	crc = hash.Sum32()
	return
}
