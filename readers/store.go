package readers

import (
	"eventdb/io"
	"eventdb/segment"
)

func ReadEvent(name string, id int) (*segment.Event, error) {

	offset, ok := findOffset(name, id)
	if ok {

	}
	evt, _ := findEvent(name, offset)

	return evt, nil
}

func findOffset(name string, id int) (uint64, bool) {

	br, err := io.NewBinaryReader(name + ".idx")

	if err != nil {
		return 0, false
	}

	defer br.Close()

	fileType, err := br.ReadHeader()

	if fileType != io.RowStoreIDXV1 {
		panic("invalid file type")
	}

	init := br.Size - 8

	return 0, true

}

func findEvent(name string, id uint64) (*segment.Event, error) {

	br, err := io.NewBinaryReader(name)

	if err != nil {
		return nil, err
	}

	defer br.Close()

	return nil, nil
}
