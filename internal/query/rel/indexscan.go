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

package rel

// TODO(sebad): Should implement Node.
type IndexScan struct {
	indexName string
	filter    *Filter
}

func NewIndexScan(name string) *IndexScan {

	s := &IndexScan{
		indexName: name,
	}

	return s
}

func (i *IndexScan) SetFilter(f *Filter) {
	i.filter = f
}

func (i *IndexScan) GetFilter() *Filter {
	return i.filter
}

func (i *IndexScan) GetIndexName() string {
	return i.indexName
}

