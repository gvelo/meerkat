package readers

import (
	"errors"
	"eventdb/io"
)

type Segment struct {
}

func ReadSegment(file string) (*Segment, error) {

	br, err := io.NewBinaryReader(file)

	if err != nil {
		return nil, err
	}

	fType, err := br.ReadHeader()

	if fType != io.SegmentInfo {
		return nil, errors.New("invalid file type")
	}
	
	indexInfo := readIndexInfo(br)

}

func readIndexInfo(br *io.BinaryReader)  {
	
}
