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

package exec

import (
	"fmt"
)

type OpNode interface {
	fmt.Stringer
	Execute(ctx Context) (Cursor, error)
	Parent() OpNode
	Children() []OpNode
	AddChild(n OpNode)
	Cancel()
}

type NodeImp struct {
	canceled bool
	parent   OpNode
	children []OpNode
}

func (p *NodeImp) Parent() OpNode {
	return p.parent
}

func (p *NodeImp) Children() []OpNode {
	return p.children
}

func (p *NodeImp) AddChild(n OpNode) {
	p.children = append(p.children, n)
}

func (p *NodeImp) Cancel() {
	p.canceled = true
}
