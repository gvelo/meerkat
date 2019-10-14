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

package logical

// Si le pongo el nombre projection el ide explota  ¯\_(ツ)_/¯
type Order struct {
	Direction string
	Field     string
}

// a logical projection.
type Projection struct {
	Indexes []string // indexes to search
	Index   string   // index selected

	Fields []string

	Limit int // limit the results

	Order []*Order // order

	RexField *RexField // field created by regex
	Span     Expression
}

func NewProjection() *Projection {
	s := &Projection{
		Indexes: make([]string, 0),
		Index:   "",
		Fields:  make([]string, 0),
	}
	return s
}

func (p *Projection) String() string {
	return "projection"
}

func (p *Projection) AddIndex(idx string) error {
	p.Indexes = append(p.Indexes, idx)
	return nil
}
