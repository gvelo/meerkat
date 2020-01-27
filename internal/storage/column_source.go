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

import "meerkat/internal/buffer"

fun

type IntColumnSource struct {
	src      *buffer.IntBuffer
	sorted   bool
	ridMap   []int
	buff     []int
	buffSize int
	pos      int
}

func (s *IntColumnSource) HasNext() bool {

}

func (s *IntColumnSource) Next() IntVector {

	if s.sorted {
		return s.nextFromSorted()
	}

	return s.nextFromUnsorted()

}

func (s *IntColumnSource) nextFromSorted() IntVector {



	return s.src.Int()[s.pos:s.buffSize]

}

func (s *IntColumnSource) nextFromUnsorted() IntVector {

}



