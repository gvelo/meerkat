package readers

import (
	"errors"
	"eventdb/io"
)

type OnDiskTrie struct {
	br         *io.BinaryReader
	rootOffset int
}

func ReadTrie(file string) (*OnDiskTrie, error) {

	br, err := io.NewBinaryReader(file)

	if err != nil {
		return nil, err
	}

	fType, err := br.ReadHeader()

	if fType != io.StringIndexV1 {
		return nil, errors.New("invalid file type")
	}

	br.Offset = br.Size - 8

	rootOffset, err := br.DecodeFixed64()

	if err != nil {
		return nil, err
	}

	return &OnDiskTrie{
		br:         br,
		rootOffset: int(rootOffset),
	}, nil

}

func (t OnDiskTrie) Lookup(term string) (int, error) {

	t.br.Offset = t.rootOffset

nodesearch:
	for i := 0; i < len(term); i++ {

		// consume Node posting
		_, err := t.br.DecodeVarint()

		if err != nil {
			return 0, err
		}

		childCount, err := t.br.DecodeVarint()

		if err != nil {
			return 0, err
		}

		for c := 0; c < int(childCount); c++ {

			key := t.br.DecodeByte()
			value, err := t.br.DecodeVarint()

			if err != nil {
				return 0, err
			}

			if term[i] == key {
				t.br.Offset = int(value)
				continue nodesearch
			}

		}

		bucketCount, err := t.br.DecodeVarint()

		if err != nil {
			return 0, err
		}

		for b := 0; b < int(bucketCount); b++ {

			value, err := t.br.DecodeStringBytes()

			if err != nil {
				return 0, err
			}

			offset, err := t.br.DecodeVarint()

			if err != nil {
				return 0, err
			}

			if value == term[i:] {
				return int(offset), nil
			}

		}

		// TODO fix the not found semantic.
		return 0, nil

	}

	offset, err := t.br.DecodeVarint()

	if err != nil {
		return 0, err
	}

	// given that we have a header in the posting list disk layout,
	// posting list offset cannot take a zero value.
	if offset != 0 {
		// final state, return the offset of the posting list
		// associated to the last node.
		return int(offset), nil
	}

	// TODO change for -1
	return 0, nil

}
