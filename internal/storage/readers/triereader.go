package readers

import (
	"errors"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment/ondsk"
)

func ReadTrie(path string) (*ondsk.BTrie, error) {

	file, err := io.MMap(path)

	if err != nil {
		return nil, err
	}

	br := file.NewBinaryReader()

	fType, err := br.ReadHeader()

	if fType != io.StringIndexV1 {
		file.UnMap()
		return nil, errors.New("invalid file type")
	}

	br.Offset = br.Size - 8

	rootOffset, err := br.ReadFixed64()

	if err != nil {
		file.UnMap()
		return nil, err
	}

	return ondsk.NewBtrie(file, int(rootOffset)), nil

}
