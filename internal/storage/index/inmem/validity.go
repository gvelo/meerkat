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

package inmem

import (
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/storage/io"
)

func NewValidityBitmapIndex(bw *io.BinaryWriter) *ValidityBitmapIndex {
	return &ValidityBitmapIndex{
		bitmap: roaring.NewBitmap(),
		bw:     bw,
	}
}

type ValidityBitmapIndex struct {
	bitmap *roaring.Bitmap
	bw     *io.BinaryWriter
}

func (v *ValidityBitmapIndex) Flush() {

	v.bitmap.RunOptimize()

	_, _ = v.bitmap.WriteTo(v.bw)

}

func (v *ValidityBitmapIndex) Cardinality() int {
	return int(v.bitmap.GetCardinality())
}

func (v *ValidityBitmapIndex) Index(rid []uint32) {
	v.bitmap.AddMany(rid)
}
