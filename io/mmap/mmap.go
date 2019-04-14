package mmap

import (
	"errors"
	"os"
)

func Map(path string) ([]byte, error) {

	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	fs, err := f.Stat()

	if err != nil {
		return nil, err
	}

	if fs.Size() == 0 {
		return nil, errors.New("trying to map an empty file")
	}

	return mmap(f, fs.Size())

}

func UnMap(ref []byte) error {
	return unmap(ref)
}
