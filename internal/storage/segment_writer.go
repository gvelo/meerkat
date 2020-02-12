// Copyright 2020 The Meerkat Authors
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

package storage

type Flushable interface {
	Flush() error
}
type IndexWriter interface {
	Flushable
}

type IntIndexWriter interface {
	IndexWriter
	Index(vector IntVector)
	Cardinality() int
}

type UintIndexWriter interface {
	IndexWriter
	Index(vector IntVector)
}

type FloatIndexWriter interface {
	IndexWriter
	Index(vector FloatVector)
}

type ByteSliceIndexWriter interface {
	IndexWriter
	Index(vector ByteSliceVector)
}

type PageIndexWriter interface {
	IndexWriter
	IndexPages(vec IntVector)
}

type PageWriter interface {
	Flushable
	//WritePage(vect ByteSliceVector)
	WritePage(page []byte, endRid uint32)
}

type ValidityIndexWriter interface {
	IndexWriter
	Index(rid []uint32)
}

type Encoder interface {
	Flushable
	FlushData() error
	Type() EncodingType
}

type IntEncoder interface {
	Encoder
	Encode(vec IntVector) error
}

type UintEncoder interface {
	Encoder
	Encode(vec IntVector) error
}

type FloatEncoder interface {
	Encoder
	Encode(vec FloatVector) error
}

type ByteSliceEncoder interface {
	Encoder
	Encode(vec ByteSliceVector) error
}

const (
	MagicNumber = "MEERKAT"
	Version     = 1
)

// TODO(gvelo) change to []byte after schema refactor.
//var TSColID = "_ts"
//
//type SegmentWriter struct {
//	table *buffer.Table
//}
//
//func (sw *SegmentWriter) WriteSegment(table *buffer.Table) error {
//
//	// write segment head
//	sw.writeTSColumn()
//
//	// crear column source y pasarselo al writer
//	// write columns
//	// write stats
//
//	return nil
//
//}
//
//func (sw *SegmentWriter) writeTSColumn() ([]int, error) {
//
//	c, ok := sw.table.Col(TSColID)
//
//	if !ok {
//		return nil, errors.New("missing TS column")
//	}
//
//	tsColumn, ok := c.(*buffer.IntBuffer)
//
//	if !ok {
//		return nil, errors.New("wrong TS column type")
//	}
//
//	sorted := sort.IntsAreSorted(tsColumn.Int())
//
//	var pos []int
//
//	if !sorted {
//		pos = SortTSColumn(tsColumn.Int())
//	}
//
//	//writeValues()
//
//	return nil,nil
//
//}
//
//func SortTSColumn(values []int) []int {
//
//	pos := make([]int, len(values))
//
//	for i := 0; i < len(pos); i++ {
//		pos[i] = i
//	}
//
//	tsSlice := &TSSlice{
//		ts:  values,
//		pos: pos,
//	}
//
//	sort.Stable(tsSlice)
//
//	return pos
//
//}
//
//type TSSlice struct {
//	ts  []int
//	pos []int
//}
//
//func (t *TSSlice) Len() int {
//	return len(t.ts)
//}
//
//func (t *TSSlice) Less(i, j int) bool {
//	return t.ts[i] < t.ts[j]
//}
//
//func (t *TSSlice) Swap(i, j int) {
//	t.ts[i], t.ts[j] = t.ts[j], t.ts[i]
//	t.pos[i], t.pos[j] = t.pos[j], t.pos[i]
//}
