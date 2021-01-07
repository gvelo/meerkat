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

package physical

import "meerkat/internal/storage/vector"

// Operator represents an Operator of a physical plan execution.
// Operators take one or more inputs and produce an output in the form
// of vectors, bitmaps or some other type.
type Operator interface {
	// Init initializes the Operator acquiring the required resources.
	// Init will call the init method on all it's input operators.
	Init()

	// Close the Operator releasing all the acquired resources.
	// Close will cascade calling the Close method on all its
	// children operators.
	Close()

	Accept(v Visitor)
}

type BatchOperator interface {
	Operator
	Next() Batch
}

type ColumnOperator interface {
	Operator
	Next() vector.Vector
}

type OutputOp interface {
	Run()
}

type Visitor interface {
	VisitPre(n Operator) Operator
	VisitPost(n Operator) Operator
}

func Walk(n Operator, v Visitor) Operator {
	n = v.VisitPre(n)
	n.Accept(v)
	return v.VisitPost(n)
}
