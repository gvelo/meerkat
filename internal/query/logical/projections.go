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
	IndexName string   // index to search
	Fields    []string // fields names to show
	Limit     int      // limit the results

	Order []*Order // order

	parent   Node
	children []Node

	RexField *RexField // field created by regex
	Span     *Exp
}

func NewProjection(name string) *Projection {
	s := &Projection{
		IndexName: name,
	}
	return s
}

func (p *Projection) String() string {
	return "projection"
}
