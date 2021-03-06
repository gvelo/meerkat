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

import (
	"meerkat/internal/storage/index"
	"meerkat/internal/storage/io"
)

type blockWriter struct {
	bw *io.BinaryWriter
	bi index.BlockIndexWriter
}

func NewBlockWriter(bw *io.BinaryWriter, bi index.BlockIndexWriter) *blockWriter {
	return &blockWriter{
		bw: bw,
		bi: bi,
	}
}

func (w *blockWriter) WriteBlock(block []byte, baseRid uint32) {
	w.bi.IndexBlock(block, w.bw.Offset(), baseRid)
	w.bw.WriteRaw(block)
}

func (w *blockWriter) Write(block []byte) {
	w.bw.WriteRaw(block)
}
