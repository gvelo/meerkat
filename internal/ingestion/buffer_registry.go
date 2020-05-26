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

package ingestion

import (
	"fmt"
	"meerkat/internal/jsoningester/ingestionpb"
)

type BufferRegistry interface {
	Add(table *ingestionpb.Table)
}

func NewBufferRegistry() BufferRegistry {
	return &bufferRegistry{}
}

type bufferRegistry struct {
}

func (b bufferRegistry) Add(table *ingestionpb.Table) {

	fmt.Println("=========================")
	fmt.Println(table.Name)
	fmt.Println(table.Columns)
	fmt.Println("=========================")

}
