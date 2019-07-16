// Copyright 2019 The Meerkat Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package writers

import (
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment/inmem"
)

type trieWriter struct {
	bw *io.BinaryWriter
}

func newTrieWriter(file string) (*trieWriter, error) {

	bw, err := io.NewBinaryWriter(file)

	if err != nil {
		return nil, err
	}

	return &trieWriter{
		bw: bw,
	}, nil

}

func (tw *trieWriter) write(trie *inmem.BTrie) error {

	err := tw.bw.WriteHeader(io.StringIndexV1)

	rootOffset, err := tw.writeNode(trie.Root)

	if err != nil {
		return err
	}

	err = tw.bw.WriteFixedUint64(uint64(rootOffset))

	if err != nil {
		return err
	}

	return nil

}

func (tw *trieWriter) writeNode(node *inmem.Node) (int, error) {

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

	// use a zero value to signal a null posting, this is valid since
	// posting never get a zero offset on disk.
	postingOffset := 0

	if node.Posting != nil {
		postingOffset = node.Posting.Offset
	}

	err := tw.bw.WriteVarUInt(postingOffset)

	if err != nil {
		return 0, err
	}

	err = tw.bw.WriteVarUInt(len(node.Children))
	if err != nil {
		return -1, err
	}

	for key, child := range node.Children {

		err = tw.bw.WriteByte(key)
		if err != nil {
			return -1, err
		}

		err = tw.bw.WriteVarUInt(child.Offset)

		if err != nil {
			return -1, err
		}

	}

	err = tw.bw.WriteVarUInt(len(node.Bucket))

	if err != nil {
		return -1, err
	}

	for _, record := range node.Bucket {

		err = tw.bw.WriteString(record.Value)

		if err != nil {
			return -1, err
		}

		err = tw.bw.WriteVarUInt(record.Posting.Offset)

		if err != nil {
			return -1, err
		}

	}

	return offset, nil

}

func (tw *trieWriter) close() error {
	return tw.bw.Close()
}

func WriteTrie(path string, trie *inmem.BTrie) error {
	tw, err := newTrieWriter(path)
	defer tw.close()
	if err != nil {
		return err
	}
	return tw.write(trie)
}
