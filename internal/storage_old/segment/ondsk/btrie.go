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

package ondsk

import "meerkat/internal/storage/io"

type BTrie struct {
	file       *io.MMFile
	rootOffset int
}

func (t *BTrie) Close() error {
	return t.file.UnMap()
}

func (t BTrie) Lookup(term string) (int, error) {

	br := t.file.NewBinaryReader()

	br.Offset = t.rootOffset

nodesearch:
	for i := 0; i < len(term); i++ {

		// consume Node posting ( only read on final state )
		_, err := br.ReadVarInt()

		if err != nil {
			return 0, err
		}

		childCount, err := br.ReadVarInt()

		if err != nil {
			return 0, err
		}

		for c := 0; c < int(childCount); c++ {

			key := br.ReadByte()
			value, err := br.ReadVarInt()

			if err != nil {
				return 0, err
			}

			if term[i] == key {
				br.Offset = value
				continue nodesearch
			}

		}

		bucketCount, err := br.ReadVarInt()

		if err != nil {
			return 0, err
		}

		for b := 0; b < bucketCount; b++ {

			value, err := br.ReadString()

			if err != nil {
				return 0, err
			}

			offset, err := br.ReadVarInt()

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

	offset, err := br.ReadVarInt()

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

func NewBtrie(file *io.MMFile, rootOffset int) *BTrie {
	return &BTrie{
		file:       file,
		rootOffset: rootOffset,
	}
}
