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

package mql_parser

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"meerkat/internal/query/logical"
	"meerkat/internal/schema"
	"testing"
	"time"
)

func TestParseQuery(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := schema.NewMockSchema(ctrl)

	// TODO(sebad): rebuild schema mock.
	s.EXPECT().
		FieldsInIndexByName(gomock.Any()).
		Return([]schema.IndexInfo{
			{
				Id:             "1",
				Name:           "coso",
				Desc:           "coso",
				Created:        time.Time{},
				Updated:        time.Time{},
				Fields:         nil,
				PartitionAlloc: schema.PartitionAlloc{},
			},
		}, nil).Times(3)

	s.EXPECT().
		IndexByName(gomock.Any()).
		Return(schema.IndexInfo{
			Id:             "1",
			Name:           "coso",
			Desc:           "coso",
			Created:        time.Time{},
			Updated:        time.Time{},
			Fields:         nil,
			PartitionAlloc: schema.PartitionAlloc{},
		}, nil).AnyTimes()

	a := assert.New(t)

	str := "index=Index campo1=100 and ( campo2=3 or campo3=\"s123\" )"

	steps, _ := Parse(s, str)

	p := steps[0].(*logical.Projection)
	a.Equal(0, len(p.Indexes))

	f := steps[1].(*logical.RootFilter).RootFilter

	a.NotNil(f)

	a.Equal(logical.AND, f.Op)
	a.Equal(false, f.Group)

	a.Equal(logical.DECIMAL, f.Left.(*logical.Filter).Right.(logical.Expression).Type())
	a.Equal("100", f.Left.(*logical.Filter).Right.(logical.Expression).Value())

	c1comp := f.Left.(*logical.Filter)
	a.Equal(logical.EQ, c1comp.Op)
	a.Equal(logical.IDENTIFIER, c1comp.Left.(logical.Expression).Type())
	a.Equal("campo1", c1comp.Left.(logical.Expression).Value())

	c1c3comp := f.Right.(*logical.Filter)
	a.Equal(logical.OR, c1c3comp.Op)
	a.Equal(true, c1c3comp.Group)

	a.Equal(logical.EQ, c1c3comp.Left.(*logical.Filter).Op)

	a.Equal(logical.EQ, c1c3comp.Right.(*logical.Filter).Op)

	a.NotNil(f.Right)
	a.NotNil(f.Op)
	a.NotNil(f.Left)

}

func TestParseQuery2(t *testing.T) {

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
			Name:      "campo1",
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
		Return([]schema.IndexInfo{ii}, nil).Times(1)

	s.EXPECT().
		IndexByName(gomock.Any()).
		Return(ii, nil).Times(1)

	a := assert.New(t)

	str := "index=Index campo1=100 | top 10 | sort campo1 desc, campo3"

	steps, _ := Parse(s, str)

	p := steps[0].(*logical.Projection)

	a.Equal(0, len(p.Indexes))
	a.Equal(10, p.Limit)

	a.Equal(2, len(p.Order))

	a.Equal("campo1", p.Order[0].Field)
	a.Equal("desc", p.Order[0].Direction)

	a.Equal("campo3", p.Order[1].Field)
	a.Equal("asc", p.Order[1].Direction)

}

func TestParseQuery3(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := schema.NewMockSchema(ctrl)

	s.EXPECT().
		FieldsInIndexByName(gomock.Any()).
		Return(make([]schema.IndexInfo, 0), nil).AnyTimes()

	a := assert.New(t)

	str := "earlier=-1d request_id=\"a37cacc3-71d5-40f0-a329-a051a3949ced\" "

	steps, _ := Parse(s, str)
	p := steps[0].(*logical.Projection)
	a.Equal(0, len(p.Indexes))

	f := steps[1].(*logical.RootFilter).RootFilter
	a.NotNil(f)

	str = "request_id=\"a37cacc3-71d5-40f0-a329-a051a3949ced\" earlier=-1d  "

	steps, _ = Parse(s, str)
	p = steps[0].(*logical.Projection)

	a.Equal(0, len(p.Indexes))
	a.NotNil(f)

}

func TestParseQuery4(t *testing.T) {

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
			Name:      "service",
			Desc:      "",
			IndexId:   "",
			FieldType: 0,
			Nullable:  false,
			Created:   time.Time{},
			Updated:   time.Time{},
		}, {
			Id:        "",
			Name:      "hbm",
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
		Return([]schema.IndexInfo{ii}, nil).Times(2)

	s.EXPECT().
		IndexByName(gomock.Any()).
		Return(ii, nil).Times(1)

	a := assert.New(t)

	str := "earlier=-1h index=Index service=hbm | bucket span=1m | stats count by _id, status"

	steps, _ := Parse(s, str)
	p := steps[0].(*logical.Projection)
	// pojection
	a.Equal(0, len(p.Indexes))

	// Filter
	a.NotNil(steps[1].(*logical.RootFilter))

	// Aggregation
	a.NotNil(steps[2].(*logical.Aggregation))

}

func TestParseQuery5(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := schema.NewMockSchema(ctrl)

	s.EXPECT().
		FieldsInIndexByName(gomock.Any()).
		Return(make([]schema.IndexInfo, 0), nil).AnyTimes()

	a := assert.New(t)

	str := "earlier=-1d | rex field=raw \"(?<time_spend>\\d{3}[0-9]+)\" " // revisar como bancarse expresiones regulares

	steps, _ := Parse(s, str)
	p := steps[0].(*logical.Projection)

	// pojection
	a.Equal(0, len(p.Indexes))

	// Filter
	a.NotNil(steps[1].(*logical.RootFilter).RootFilter)

}

func TestParseQuery6(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := schema.NewMockSchema(ctrl)

	s.EXPECT().
		IndexByName(gomock.Any()).
		Return(schema.IndexInfo{
			Id:      "ii",
			Name:    "ii",
			Desc:    "index ii",
			Created: time.Time{},
			Updated: time.Time{},
			Fields: []schema.Field{{
				Id:        "",
				Name:      "campo1",
				Desc:      "",
				IndexId:   "",
				FieldType: 0,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			}, {
				Id:        "",
				Name:      "campo2",
				Desc:      "",
				IndexId:   "",
				FieldType: 0,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			}, {
				Id:        "",
				Name:      "campo3",
				Desc:      "",
				IndexId:   "",
				FieldType: 0,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			},
			PartitionAlloc: schema.PartitionAlloc{},
		}, nil).Times(0)

	a := assert.New(t)

	str := "earlier=-1d | fields -campo1 -campo2 "

	_, err := Parse(s, str)
	if err == nil {
		a.Fail("Error not thrown")
	}

}

func TestParseQuery7(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := schema.NewMockSchema(ctrl)

	s.EXPECT().
		IndexByName(gomock.Any()).
		Return(schema.IndexInfo{
			Id:      "ii",
			Name:    "ii",
			Desc:    "index ii",
			Created: time.Time{},
			Updated: time.Time{},
			Fields: []schema.Field{{
				Id:        "",
				Name:      "campo2",
				Desc:      "",
				IndexId:   "",
				FieldType: 0,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			}, {
				Id:        "",
				Name:      "campo3",
				Desc:      "",
				IndexId:   "",
				FieldType: 0,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			},
			PartitionAlloc: schema.PartitionAlloc{},
		}, nil).Times(1)

	a := assert.New(t)

	str := "earlier=-1d index=ii | fields -campo1 -campo2 "

	_, err := Parse(s, str)
	if err != nil {
		a.Fail("Error thrown")
	}

}
