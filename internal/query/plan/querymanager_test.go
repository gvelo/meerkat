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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"meerkat/internal/schema"
	"testing"
	"time"
)

func Test_Query_Fields_bad_type(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := schema.NewMockSchema(ctrl)

	s.EXPECT().
		FieldsInIndexByName(gomock.Eq("f1")).
		Return(make([]schema.IndexInfo, 0), &schema.NotFoundError{Err: "No se encontro el campo."})

	assert := assert.New(t)

	qm, err := NewQueryManager(s, nil, nil, nil)
	assert.NoError(err)

	_, err = qm.Query("earlier=1d f1=\"campo1\" ")
	assert.NotNil(err)
}

func Test_Query_FieldsOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assert := assert.New(t)

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

	qm, err := NewQueryManager(s, nil, nil, nil)
	assert.NoError(err)

	_, err = qm.Query("earlier=1d f1=\"campo1\" f2=12 or f1=f2 ")
	assert.Nil(err)
}
