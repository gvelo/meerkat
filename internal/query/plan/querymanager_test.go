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
		FieldByName(gomock.Eq("f1")).
		Return(schema.Field{}, &schema.NotFoundError{Err: "No se encontro el campo."})

	assert := assert.New(t)

	qm, err := NewQueryManager(s)
	assert.NoError(err)

	_, err = qm.Query("earlier=1d f1=\"campo1\" ")
	assert.NotNil(err)
}

func Test_Query_FieldsOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assert := assert.New(t)

	s := schema.NewMockSchema(ctrl)

	s.EXPECT().
		FieldByName(gomock.Eq("f1")).
		Return(schema.Field{
			Id:        "f1",
			Name:      "f1",
			Desc:      "f1",
			IndexId:   "f1",
			FieldType: 0,
			Nullable:  false,
			Created:   time.Time{},
			Updated:   time.Time{},
		}, nil).Times(2)

	s.EXPECT().
		FieldByName(gomock.Eq("f2")).
		Return(schema.Field{
			Id:        "f2",
			Name:      "f2",
			Desc:      "f2",
			IndexId:   "f2",
			FieldType: 0,
			Nullable:  false,
			Created:   time.Time{},
			Updated:   time.Time{},
		}, nil).Times(2)

	qm, err := NewQueryManager(s)
	assert.NoError(err)

	_, err = qm.Query("earlier=1d f1=\"campo1\" f2=12 or f1=f2 ")
	assert.Nil(err)
}
