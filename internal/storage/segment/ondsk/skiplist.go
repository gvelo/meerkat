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
)

// Devuelve el offset de la pagina.
type SkipList struct {
	Br         *io.BinaryReader
	RootOffset int
	Lvl        int
	Interface  SkipListInterface
}

type SkipListInterface interface {
	ReadValue(sl *SkipList) interface{}
	Compare(interface{}, interface{}) int
}

type IntInterface struct {
}

func (i IntInterface) ReadValue(sl *SkipList) interface{} {
	n, err := sl.Br.ReadVarInt()
	if err != nil {
		panic(err)
	}
	return n
}

func (i IntInterface) Compare(a interface{}, b interface{}) int {
	if a.(int) < b.(int) {
		return -1
	}
	if a.(int) > b.(int) {
		return 1
	}
	return 0
}

type UInt32Interface struct {
}

func (i UInt32Interface) ReadValue(sl *SkipList) interface{} {
	n, err := sl.Br.ReadVarInt()
	if err != nil {
		panic(err)
	}
	return n
}

func (i UInt32Interface) Compare(a interface{}, b interface{}) int {
	if uint32(a.(int)) < b.(uint32) {
		return -1
	}
	if uint32(a.(int)) > b.(uint32) {
		return 1
	}
	return 0
}

type FloatInterface struct {
}

func (i FloatInterface) ReadValue(sl *SkipList) interface{} {
	n, err := sl.Br.ReadFixed64()
	r := math.Float64frombits(n)
	if err != nil {
		panic(err)
	}
	return r
}

func (i FloatInterface) Compare(a interface{}, b interface{}) int {
	if a.(float64) < b.(float64) {
		return -1
	}
	if a.(float64) > b.(float64) {
		return 1
	}
	return 0
}

func (sl *SkipList) Lookup(id interface{}) (uint64, error) {
	_, o, ok, _ := sl.findOffsetSkipList(id)
	if ok {
		return o, nil
	}
	return 0, nil
}

func (sl *SkipList) findOffsetSkipList(id interface{}) (interface{}, uint64, bool, error) {

	r, start, err := sl.readSkipList(int(sl.RootOffset), sl.Lvl, id)
	return r, start, true, err

}

func (sl *SkipList) readSkipList(offset int, lvl int, id interface{}) (interface{}, uint64, error) {

	sl.Br.Offset = offset

	// search this lvl
	for i := 0; i < int(math.MaxUint32); i++ {

		if lvl == 0 {

			k := sl.Interface.ReadValue(sl)
			kOffset, _ := sl.Br.ReadVarint64()

			if sl.Interface.Compare(k, id) == 0 {
				return k, kOffset, nil
			}

			if sl.Interface.Compare(k, id) > 0 {
				// not found
				return k, kOffset, nil
			}
		} else {
			sl.Br.Offset = offset
			k := sl.Interface.ReadValue(sl)
			kOffset, _ := sl.Br.ReadVarint64()
			next := sl.Br.Offset
			kn := sl.Interface.ReadValue(sl)
			sl.Br.ReadVarint64()

			if sl.Interface.Compare(k, id) == 0 {
				return sl.readSkipList(int(kOffset), lvl-1, id)
			}

			if sl.Interface.Compare(kn, id) > 0 {
				// done, not found
				if lvl == 0 {
					return 0, 0, nil
				}
				return sl.readSkipList(int(kOffset), lvl-1, id)
			}
			offset = next
		}

	}
	return 0, 0, nil

}
