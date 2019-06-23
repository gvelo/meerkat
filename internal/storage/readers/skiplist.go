package readers

import (
	"errors"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment/ondsk"
)

func ReadSkipList(path string) (*ondsk.OnDiskSkipList, error) {

	file, err := io.MMap(path)

	if err != nil {
		return nil, err
	}

	br := file.NewBinaryReader()

	fType, err := br.ReadHeader()

	if fType != io.SkipListV1 {
		return nil, errors.New("invalid file type")
	}

	br.Offset = br.Size - 16
	rootOffset, _ := br.ReadFixed64()
	lvl, _ := br.ReadFixed64()

	if err != nil {
		return nil, err
	}

	return &ondsk.OnDiskSkipList{
		Br:         br,
		RootOffset: int(rootOffset),
		Lvl:        int(lvl),
	}, nil
}
