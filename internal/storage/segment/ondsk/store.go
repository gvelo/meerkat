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

import (
	"math"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment"
)

type OnDiskColumn struct {
	Fi *segment.FieldInfo
	Br *io.BinaryReader
}

type OnDiskColumnIdx struct {
	Br         *io.BinaryReader
	RootOffset int
	Ixl        int
	Lvl        int
	Fi         *segment.FieldInfo
}

func (sl *OnDiskColumnIdx) Lookup(id int) (int, error) {
	offset, _, err := sl.findOffset(id)
	return offset, err
}

func (sl *OnDiskColumnIdx) findOffset(id int) (int, int, error) {

	r, start := sl.processLevel(int(sl.RootOffset), sl.Lvl, id, 0)

	return r, start, nil

}

func (sl *OnDiskColumnIdx) processLevel(offset int, lvl int, id int, start int) (int, int) {
	// if it is the 1st lvl & the offsets are less than
	// the Ixl then return the offset 0 to search from
	if lvl == 0 {
		return offset, start
	}
	sl.Br.Offset = offset

	// search this Lvl
	for i := 0; i < int(math.MaxUint32); i++ {

		sId1, _ := sl.Br.ReadVarUInt()
		oId1, _ := sl.Br.ReadVarUInt()

		sId2, _ := sl.Br.ReadVarUInt()

		if int(sId2) > id {
			return sl.processLevel(int(oId1), lvl-1, id, int(sId1))
		}

	}
	return 0, 0
}
