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

package index

import (
	"meerkat/internal/storage/colval"
)

type IndexWriter interface {
	Flush()
	Cardinality() int
}

type IntIndexWriter interface {
	IndexWriter
	Index(v colval.IntColValues)
}

type UintIndexWriter interface {
	IndexWriter
	Index(v colval.UintColValues)
}

type FloatIndexWriter interface {
	IndexWriter
	Index(v colval.FloatColValues)
}

type ByteSliceIndexWriter interface {
	IndexWriter
	Index(v colval.ByteSliceColValues)
}

type BoolIndexWriter interface {
	IndexWriter
	Index(v colval.BoolColValues)
}

type BlockIndexWriter interface {
	Flush()
	IndexBlock(block []byte, baseRID uint32)
}

type ValidityIndexWriter interface {
	IndexWriter
	Index(rid []uint32)
}
