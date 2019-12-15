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
	"meerkat/internal/schema"
	"time"
)

type Context interface {
	fmt.Stringer
	Schema() schema.Schema
	SetSpan(sp time.Duration)
	Span() time.Duration
}

type QryContext struct {
	schema schema.Schema
	span   time.Duration
}

func (ctx *QryContext) String() string {
	return "QryContext"
}

func NewQueryContext(s schema.Schema) Context {
	return &QryContext{
		schema: s,
	}
}

func (ctx *QryContext) Schema() schema.Schema {
	return ctx.schema
}

func (ctx *QryContext) SetSpan(sp time.Duration) {
	ctx.span = sp
}

func (ctx *QryContext) Span() time.Duration {
	return ctx.span
}
