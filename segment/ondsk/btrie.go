package ondsk

import "eventdb/io"

type OnDiskTrie struct {
	Br         *io.BinaryReader
	RootOffset uint64
}

func (t OnDiskTrie) Lookup(term string) (uint64, error) {

	t.Br.Offset = int64(t.RootOffset)

nodesearch:
	for i := 0; i < len(term); i++ {

		// consume Node posting
		_, err := t.Br.DecodeVarint()

		if err != nil {
			return 0, err
		}

		childCount, err := t.Br.DecodeVarint()

		if err != nil {
			return 0, err
		}

		for c := 0; c < int(childCount); c++ {

			key := t.Br.DecodeByte()
			value, err := t.Br.DecodeVarint()

			if err != nil {
				return 0, err
			}

			if term[i] == key {
				t.Br.Offset = int64(value)
				continue nodesearch
			}

		}

		bucketCount, err := t.Br.DecodeVarint()

		if err != nil {
			return 0, err
		}

		for b := 0; b < int(bucketCount); b++ {

			value, err := t.Br.DecodeStringBytes()

			if err != nil {
				return 0, err
			}

			offset, err := t.Br.DecodeVarint()

			if err != nil {
				return 0, err
			}

			if value == term[i:] {
				return offset, nil
			}

		}

		// TODO fix the not found semantic.
		return 0, nil

	}

	offset, err := t.Br.DecodeVarint()

	if err != nil {
		return 0, err
	}

	// given that we have a header in the posting list disk layout,
	// posting list offset cannot take a zero value.
	if offset != 0 {
		// final state, return the offset of the posting list
		// associated to the last node.
		return offset, nil
	}

	// TODO change for -1
	return 0, nil

}
