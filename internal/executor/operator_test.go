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

package executor

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage"
	"testing"
)

type BitmapIndexOperator struct {
	i     storage.IntIndex
	op    Operation
	value interface{}
}

func (b BitmapIndexOperator) Init() {
	panic("implement me")
}

func (b BitmapIndexOperator) Destroy() {
	panic("implement me")
}

func (b BitmapIndexOperator) Next() *roaring.Bitmap {
	panic("implement me")
}

// Select * From Index Where c1 = 12 and ts Between  F1 AND F2
// C1, F1, F2 indexed
func TestBinaryBitmapOperator(t *testing.T) {
	a := assert.New(t)
	var ts storage.IntColumn
	var c1 storage.IntColumn

	op1 := NewIntIndexScanOperator(lt, 1, ts) // ts > 1
	op2 := NewIntIndexScanOperator(gt, 2, ts) // ts < 2

	op3 := NewBinaryBitmapOperator(and, op1, op2) // ts > 1 AND ts < 2
	op4 := NewIntIndexScanOperator(and, 12, c1)   // C1 == 12

	op5 := NewBinaryBitmapOperator(gt, op3, op4) // ts > 1 AND ts < 2 AND C1 == 12

	op6 := NewReaderOperator(op5) // definir columna o mejor este tiene que saber que hacer en todos los campos
	// op5 := Decompress

}

func NewReaderOperator(child BitmapOperator) *ReaderOperator {
	return &ReaderOperator{child}
}

type ReaderOperator struct {
	child BitmapOperator // (Positions to review)
}

func (r *ReaderOperator) Init() {

}

func (r *ReaderOperator) Destroy() {

}

func (r *ReaderOperator) Next() storage.Vector {
	return nil
}
