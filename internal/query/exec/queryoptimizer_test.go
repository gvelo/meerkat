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
	"github.com/golang/mock/gomock"
	"meerkat/internal/query/logical"
	"meerkat/internal/schema"
	"testing"
	"time"
)

func Test_Optimize_Fields_bad_type(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := schema.NewMockSchema(ctrl)

	ii := schema.IndexInfo{
		Id:      "1",
		Name:    "Index",
		Desc:    "coso",
		Created: time.Time{},
		Updated: time.Time{},
		Fields: []schema.Field{{
			Id:        "",
			Name:      "f1",
			Desc:      "",
			IndexId:   "",
			FieldType: 0,
			Nullable:  false,
			Created:   time.Time{},
			Updated:   time.Time{},
		}, {
			Id:        "",
			Name:      "f2",
			Desc:      "",
			IndexId:   "",
			FieldType: 0,
			Nullable:  false,
			Created:   time.Time{},
			Updated:   time.Time{},
		}},
		PartitionAlloc: schema.PartitionAlloc{},
	}

	s.EXPECT().
		FieldsInIndexByName(gomock.Any()).
		Return([]schema.IndexInfo{ii}, nil).Times(4)

	p := logical.NewProjection()
	p.Fields = []string{"f1", "f2"}
	p.Index = "Index"
	p.Limit = 10

	ctx := NewQueryContext(s)

	o := NewMeerkatOptimizer()

	o.transformProjection(ctx, p)

}
