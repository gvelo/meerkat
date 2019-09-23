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

package plan

import (
	"meerkat/internal/query/logical"
	"meerkat/internal/query/mql_parser"
	"testing"
)

func buildIndexScan() *logical.Projection {

	sql := "indexname=name  campo1=12 and  ( campo2>12  or campo1=2) "
	return mql_parser.Parse(sql)

}

func TestMeerkatOptimizer_OptimizeQuery(t *testing.T) {

}
