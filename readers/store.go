package readers

import (
	"errors"
	"eventdb/io"
	"eventdb/segment"
	"eventdb/segment/ondsk"
	"strings"
)

func ReadStore(path string, ii *segment.IndexInfo, ixl int) (*ondsk.OnDiskStore, error) {

	file, err := io.MMap(path)

	if err != nil {
		return nil, err
	}

	br := file.NewBinaryReader()

	fType, err := br.ReadHeader()

	if fType != io.RowStoreIDXV1 {
		return nil, errors.New("invalid file type")
	}

	br.Offset = br.Size - 16
	rootOffset, _ := br.ReadFixed64()
	lvl, _ := br.ReadFixed64()

	if err != nil {
		return nil, err
	}

	cols := make([]*ondsk.OnDiskColumn, 0)

	for _, fi := range ii.Fields {
		col, _ := ReadColumn(path, fi)
		cols = append(cols, col)
	}

	return &ondsk.OnDiskStore{
		Br:         br,
		RootOffset: int(rootOffset),
		IndexInfo:  ii,
		Lvl:        int(lvl),
		Ixl:        ixl,
		Columns:    cols,
	}, nil
}

func ReadColumn(path string, info *segment.FieldInfo) (*ondsk.OnDiskColumn, error) {

	n := strings.Replace(path, idxExt, "."+info.Name+binExt, 1)

	file, err := io.MMap(n)

	if err != nil {
		return nil, err
	}

	br := file.NewBinaryReader()

	fileType, _ := br.ReadHeader()

	if fileType != io.RowStoreV1 {
		panic("invalid file type")
	}

	return &ondsk.OnDiskColumn{Br: br, Fi: info}, nil

}
