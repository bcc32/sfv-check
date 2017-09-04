// +build linux darwin

package sfv

import (
	"hash/crc32"
	. "syscall"
)

func Crc32File(filename string) (uint32, error) {
	fd, err := Open(filename, O_RDONLY, 0)
	if err != nil {
		return 0, errFileOpen{err, filename}
	}

	var stat Stat_t
	err = Fstat(fd, &stat)
	if err != nil {
		return 0, errFileOpen{err, filename}
	}

	if stat.Size == 0 {
		return 0, nil
	}

	buf, err := Mmap(fd, 0, int(stat.Size),
		PROT_READ, MAP_SHARED)
	if err != nil {
		return 0, errFileOpen{err, filename}
	}

	hash := crc32.NewIEEE()
	hash.Write(buf)

	return hash.Sum32(), nil
}
