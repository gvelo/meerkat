// +build darwin linux

package mmap

import (
	"os"
	"syscall"
)

func mmap(f *os.File, size int64) ([]byte, error) {

	b, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)

	// TODO add kernel advicing.

	if err != nil {
		return nil, err
	}

	return b, nil

}

func unmap(ref []byte) error {
	return syscall.Munmap(ref)
}
