// Copyright 2019 The Meerkat Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this path except in compliance with the License.
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
	"log"
	"meerkat/internal/config"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment/inmem"
)

func WriteStoreIdx(name string, offsets []*inmem.Page) error {

	bw, err := io.NewBinaryWriter(name)

	if err != nil {
		return err
	}

	defer bw.Close()
	err = bw.WriteHeader(io.RowStoreIDXV1)

	err, lvl, lvlOffset := nil, 1, bw.Offset
	lvl0Offsets, _ := writeLevel0(bw, offsets)
	if config.SkipLevelSize < len(lvl0Offsets) {
		err, lvl, lvlOffset = processLevel(bw, lvl0Offsets, lvl, 0, 0)
		if err != nil {
			panic(err)
		}
	}

	bw.WriteFixedUint64(uint64(lvlOffset))
	bw.WriteFixedUint64(uint64(lvl))

	return nil
}

func writeLevel0(bw *io.BinaryWriter, offsets []*inmem.Page) ([]*inmem.Page, error) {
	o := make([]*inmem.Page, 0)
	for i := 0; i < len(offsets); i++ {
		// starting id , idx offset
		o = append(o, &inmem.Page{StartID: offsets[i].StartID, Offset: bw.Offset})
		bw.WriteVarUint64(uint64(offsets[i].StartID))
		bw.WriteVarUint64(uint64(offsets[i].Offset))
	}
	log.Printf("Write Lvl 0 %v", offsets)
	return o, nil
}

func processLevel(bw *io.BinaryWriter, offsets []*inmem.Page, lvl int, ts uint64, lastOffset int) (error, int, int) {

	offset := int(bw.Offset)
	if len(offsets) <= 1 {
		return nil, lvl - 1, lastOffset
	}

	nl := make([]*inmem.Page, 0) // offsets storeFile

	for i := 0; i < len(offsets); i++ {

		if i%config.SkipLevelSize == 0 {
			nl = append(nl, &inmem.Page{StartID: offsets[i].StartID, Offset: bw.Offset})
		}
		bw.WriteVarUint64(uint64(offsets[i].StartID))
		bw.WriteVarUint64(uint64(offsets[i].Offset))

	}
	return processLevel(bw, nl, lvl+1, ts, offset)
}
