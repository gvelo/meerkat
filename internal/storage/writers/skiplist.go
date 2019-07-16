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
	"math"
	"meerkat/internal/config"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment/inmem"
)

func WriteSkipList(name string, sl *inmem.SkipList) error {

	bw, err := io.NewBinaryWriter(name)

	if err != nil {
		return err
	}

	defer bw.Close()

	err = bw.WriteHeader(io.SkipListV1)

	if err != nil {
		return err
	}

	offsets := make([]uint64, 0)
	keys := make([]float64, 0)

	it := sl.NewIterator(0)
	for it.HasNext() {
		// TODO: FIX copio las claves, puede ocupar mucho...
		var key = it.Key()
		switch key.(type) {
		case uint64:
			key = float64(key.(uint64))
		case int:
			key = float64(key.(int))
		case float64:
			key = key.(float64)
		}

		keys = append(keys, key.(float64))
		offsets = append(offsets, uint64(it.Next().UserData.(*inmem.PostingList).Offset))
	}

	writeSkipIdx(bw, keys, offsets)

	return nil
}

func writeSkipIdx(bw *io.BinaryWriter, keys []float64, offsets []uint64) error {

	err, l, lvlOffset := processSkip(bw, keys, offsets, 0, int(bw.Offset))
	if err != nil {
		panic(err)
	}
	bw.WriteFixedUint64(uint64(lvlOffset))
	bw.WriteFixedUint64(uint64(l))

	return nil
}

func processSkip(bw *io.BinaryWriter, keys []float64, offsets []uint64, lvl int, lastOffset int) (error, int, int) {

	offset := bw.Offset
	if len(offsets) <= 2 {
		return nil, lvl - 1, lastOffset
	}

	nl := make([]uint64, 0)
	nk := make([]float64, 0)

	for i := 0; i < len(offsets); i++ {
		o := bw.Offset
		bw.WriteFixedUint64(math.Float64bits(keys[i]))
		bw.WriteVarUint64(offsets[i])
		if i%config.SkipLevelSize == 0 {
			nk = append(nk, keys[i])
			nl = append(nl, uint64(o))
		}

	}

	bw.WriteFixedUint64(math.Float64bits(math.MaxFloat64))

	return processSkip(bw, nk, nl, lvl+1, offset)
}
