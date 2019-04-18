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

func (tw *TrieWriter) Write(trie *inmem.BTrie) (int64, error) {
	return tw.writeNode(trie.Root)
}

func (tw *TrieWriter) writeNode(node *inmem.Node) (int64, error) {

	for _, node := range node.Children {
		offset, err := tw.writeNode(node)
		if err != nil {
			return -1, err
		}
		node.Offset = offset
	}

	offset := tw.bw.Offset

	err := tw.bw.WriteEncodedVarint(uint64(len(node.Children)))

	if err != nil {
		return -1, err
	}

	for key, node := range node.Children {

		err = tw.bw.WriteByte(key)
		if err != nil {
			return -1, err
		}

		// TODO please, check this cast.
		err = tw.bw.WriteEncodedVarint(uint64(node.Offset))

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
