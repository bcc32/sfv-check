// +build linux darwin

package sfv

import (
	"hash/crc32"
	. "syscall"
)

func Crc32File(filename string) (crc uint32, error error) {
	fd, err := Open(filename, O_RDONLY, 0)
	if err != nil {
		error = err
		return
	}
	defer func() {
		err := Close(fd)
		if err != nil {
			error = err
		}
	}()

	var stat Stat_t
	err = Fstat(fd, &stat)
	if err != nil {
		error = err
		return
	}

	if stat.Size == 0 {
		return
	}

	buf, err := Mmap(fd, 0, int(stat.Size),
		PROT_READ, MAP_SHARED)
	if err != nil {
		error = err
		return
	}
	defer func() {
		err := Munmap(buf)
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
