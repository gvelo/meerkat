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

package inmem

import "fmt"

// Type represent the type of a field.
type Encoding uint

const (
	RLE Encoding = iota
	Simple8B
	DoubleDelta
	Raw
	ZigZag
	Snappy
)

type Page struct {
	Enc         Encoding
	StartID     int
	PayloadSize int
	Total       int
}

type PageDescriptor struct {
	StartID int
	Offset  int
}

func (pd *PageDescriptor) String() string {
	return fmt.Sprintf("{ S: %d, O :%d  }", pd.StartID, pd.Offset)
}
