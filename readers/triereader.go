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

	rootOffset, err := br.ReadFixed64()

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

		// consume Node posting ( only read on final state )
		_, err := t.br.ReadVarInt()

		if err != nil {
			return 0, err
		}

		childCount, err := t.br.ReadVarInt()

		if err != nil {
			return 0, err
		}

		for c := 0; c < int(childCount); c++ {

			key := t.br.ReadByte()
			value, err := t.br.ReadVarInt()

			if err != nil {
				return 0, err
			}

			if term[i] == key {
				t.br.Offset = value
				continue nodesearch
			}

		}

		bucketCount, err := t.br.ReadVarInt()

		if err != nil {
			return 0, err
		}

		for b := 0; b < bucketCount; b++ {

			value, err := t.br.ReadString()

			if err != nil {
				return 0, err
			}

			offset, err := t.br.ReadVarInt()

			if err != nil {
				return 0, err
			}

			if value == term[i:] {
				return offset, nil
			}

		}

		// TODO fix the not found semantic.
		return -1, nil

	}

	offset, err := t.br.ReadVarInt()

	if err != nil {
		return 0, err
	}

	// zero means null posting
	if offset != 0 {
		// final state, return the offset of the posting list
		// associated to the last node.
		return offset, nil
	}

	return -1, nil

}
