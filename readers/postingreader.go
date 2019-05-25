package readers

import (
	"errors"
	"eventdb/io"
	"eventdb/segment/ondsk"
)

func ReadPostingStore(path string) (*ondsk.PostingStore, error) {

	file, err := io.MMap(path)

	if err != nil {
		return nil, err
	}

	br := file.NewBinaryReader()

	fileType, err := br.ReadHeader()

	if fileType != io.PostingListV1 {
		file.UnMap()
		return nil, errors.New("invalid file type")
	}

	ps := ondsk.NewPostingStore(file)

	return ps, nil

}
