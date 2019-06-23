package io

import "meerkat/internal/storage/io/mmap"

type MMFile struct {
	Bytes []byte
}

func MMap(path string) (*MMFile, error) {

	b, err := mmap.Map(path)

	if err != nil {
		return nil, err
	}

	return &MMFile{
		Bytes: b,
	}, nil

}

func (f *MMFile) UnMap() error {
	return mmap.UnMap(f.Bytes)
}

func (f *MMFile) NewBinaryReader() *BinaryReader {
	return NewBinaryReader(f.Bytes)
}
