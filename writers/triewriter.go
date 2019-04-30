package writers

import (
	"eventdb/io"
	"eventdb/segment/inmem"
)

type TrieWriter struct {
	bw *io.BinaryWriter
}

func NewTrieWriter(file string) (*TrieWriter, error) {

	bw, err := io.NewBinaryWriter(file)

	if err != nil {
		return nil, err
	}

	return &TrieWriter{
		bw: bw,
	}, nil

}

func (tw *TrieWriter) Write(trie *inmem.BTrie) error {

	err := tw.bw.WriteHeader(io.StringIndexV1)

	rootOffset, err := tw.writeNode(trie.Root)

	if err != nil {
		return err
	}

	err = tw.bw.WriteEncodedFixed64(uint64(rootOffset))

	if err != nil {
		return err
	}

	return nil

}

func (tw *TrieWriter) writeNode(node *inmem.Node) (int, error) {

	// Write children first
	for _, child := range node.Children {
		offset, err := tw.writeNode(child)
		if err != nil {
			return -1, err
		}
		child.Offset = offset
	}

	// the node  start offset.
	offset := tw.bw.Offset

	// use a zero value to signal a null posting, this is possible since
	// posting never get a zero offset on disk.
	postingOffset := uint64(0)

	if node.Posting != nil {
		postingOffset = uint64(node.Posting.Offset)
	}

	err := tw.bw.WriteEncodedVarint(postingOffset)

	if err != nil {
		return 0, err
	}

	err = tw.bw.WriteEncodedVarint(uint64(len(node.Children)))
	if err != nil {
		return -1, err
	}

	for key, child := range node.Children {

		err = tw.bw.WriteByte(key)
		if err != nil {
			return -1, err
		}

		// TODO check this cast.
		err = tw.bw.WriteEncodedVarint(uint64(child.Offset))

		if err != nil {
			return -1, err
		}

	}

	err = tw.bw.WriteEncodedVarint(uint64(len(node.Bucket)))

	if err != nil {
		return -1, err
	}

	for _, record := range node.Bucket {

		err = tw.bw.WriteEncodedStringBytes(record.Value)

		if err != nil {
			return -1, err
		}

		err = tw.bw.WriteEncodedVarint(uint64(record.Posting.Offset))

		if err != nil {
			return -1, err
		}

	}

	return offset, nil

}

func (tw *TrieWriter) Close() error {
	return tw.bw.Close()
}
