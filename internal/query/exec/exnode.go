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

type ExNode interface {
	fmt.Stringer
	Execute(ctx Context) (Cursor, error)
	Parent() ExNode
	AddChild(n ExNode)
	Cancel()
}

type MultiNode interface {
	ExNode
	Children() []ExNode
}

type SingleNode interface {
	ExNode
	Child() ExNode
}

type SingleNodeImpl struct {
	canceled bool
	parent   ExNode
	child    ExNode
}

func (p *SingleNodeImpl) Parent() ExNode {
	return p.parent
}

func (p *SingleNodeImpl) Child() ExNode {
	return p.child
}

func (p *SingleNodeImpl) AddChild(n ExNode) {
	p.child = n
}

func (p *SingleNodeImpl) Cancel() {
	p.canceled = true
}

type MultiNodeImpl struct {
	canceled bool
	parent   ExNode
	children []ExNode
}

func (p *MultiNodeImpl) Parent() ExNode {
	return p.parent
}

func (p *MultiNodeImpl) Child() []ExNode {
	return p.children
}

func (p *MultiNodeImpl) AddChild(n ExNode) {
	p.children = append(p.children, n)
}

func (p *MultiNodeImpl) Cancel() {
	p.canceled = true
}
