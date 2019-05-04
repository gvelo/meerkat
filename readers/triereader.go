package readers

import (
	"errors"
	"eventdb/io"
	"eventdb/segment/ondsk"
)

func ReadTrie(file string) (*ondsk.OnDiskTrie, error) {

	br, err := io.NewBinaryReader(file)

	if err != nil {
		return nil, err
	}

	fType, err := br.ReadHeader()

	if fType != io.StringIndexV1 {
		return nil, errors.New("invalid file type")
	}

	br.Offset = int64(br.Size) - 8

	rootOffset, err := br.DecodeFixed64()

	if err != nil {
		return nil, err
	}

	return &ondsk.OnDiskTrie{
		Br:         br,
		RootOffset: rootOffset,
	}, nil

}
