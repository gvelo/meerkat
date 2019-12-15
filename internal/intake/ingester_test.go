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

package intake

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"meerkat/internal/schema"
	"strconv"
	"testing"
	"time"
)

func createIndex(nullable bool) schema.IndexInfo {

	index := schema.IndexInfo{
		Id:             "test-index",
		Name:           "test-index",
		Desc:           "test-index",
		Created:        time.Time{},
		Updated:        time.Time{},
		PartitionAlloc: schema.PartitionAlloc{},
		Fields: []schema.Field{
			{
				Id:        "intFieldId",
				Name:      "intField",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_INT,
				Nullable:  nullable,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "uintFieldId",
				Name:      "uintField",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_UINT,
				Nullable:  nullable,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "floatFieldId",
				Name:      "floatField",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_FLOAT,
				Nullable:  nullable,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "stringFieldId",
				Name:      "stringField",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_STRING,
				Nullable:  nullable,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "textFieldId",
				Name:      "textField",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_TEXT,
				Nullable:  nullable,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "tsFieldId",
				Name:      "tsField",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_TIMESTAMP,
				Nullable:  nullable,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "uuidFieldId",
				Name:      "uuidField",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_UUID,
				Nullable:  nullable,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
		},
	}

	return index

}

func Test(t *testing.T) {

	type Col struct {
		name  string
		nulls uint64
		err   bool
	}

	type TestCase struct {
		name        string
		json        string
		index       schema.IndexInfo
		invalidJson bool
		len         int
		nulls       uint64
		col         *Col
	}

	tests := []TestCase{
		// INT fields

		{
			name:        "test valid rows",
			json:        createJson("", ""),
			index:       createIndex(false),
			len:         6,
			invalidJson: false,
		},
		{
			name:        "test invalid int",
			json:        createJson("intField", "invalid_int"),
			index:       createIndex(false),
			len:         5,
			invalidJson: false,
			col: &Col{
				name:  "intField",
				nulls: 0,
				err:   true,
			},
		},
		{
			name:        "test invalid int (empty)",
			json:        createJson("intField", ""),
			index:       createIndex(false),
			len:         5,
			invalidJson: false,
			col: &Col{
				name:  "intField",
				nulls: 0,
				err:   true,
			},
		},
		{
			name:        "test invalid int (nil)",
			json:        createJson("intField", nil),
			index:       createIndex(false),
			len:         5,
			invalidJson: false,
			col: &Col{
				name:  "intField",
				nulls: 0,
				err:   true,
			},
		},
		{
			name:        "test valid int (string)",
			json:        createJson("intField", "10"),
			index:       createIndex(false),
			len:         6,
			invalidJson: false,
			col: &Col{
				name:  "intField",
				nulls: 0,
				err:   false,
			},
		},

		// UINT fields

		{
			name:        "test invalid uint (negative)",
			json:        createJson("uintField", -1),
			index:       createIndex(false),
			len:         5,
			invalidJson: false,
			col: &Col{
				name:  "uintField",
				nulls: 0,
				err:   true,
			},
		},
		{
			name:        "test invalid uint (negative string)",
			json:        createJson("uintField", "-1"),
			index:       createIndex(false),
			len:         5,
			invalidJson: false,
			col: &Col{
				name:  "uintField",
				nulls: 0,
				err:   true,
			},
		},
		{
			name:        "test invalid uint (invalid string)",
			json:        createJson("uintField", "invalid uint"),
			index:       createIndex(false),
			len:         5,
			invalidJson: false,
			col: &Col{
				name:  "uintField",
				nulls: 0,
				err:   true,
			},
		},
		{
			name:        "test valid uint (string)",
			json:        createJson("uintField", "10"),
			index:       createIndex(false),
			len:         6,
			invalidJson: false,
			col: &Col{
				name:  "uintField",
				nulls: 0,
				err:   false,
			},
		},

		// float fields

		{
			name:        "test invalid float ( invalid string )",
			json:        createJson("floatField", "invalid_float"),
			index:       createIndex(false),
			len:         5,
			invalidJson: false,
			col: &Col{
				name:  "floatField",
				nulls: 0,
				err:   true,
			},
		},
		{
			name:        "test invalid float (empty)",
			json:        createJson("floatField", ""),
			index:       createIndex(false),
			len:         5,
			invalidJson: false,
			col: &Col{
				name:  "floatField",
				nulls: 0,
				err:   true,
			},
		},
		{
			name:        "test invalid float (nil)",
			json:        createJson("floatField", nil),
			index:       createIndex(false),
			len:         5,
			invalidJson: false,
			col: &Col{
				name:  "floatField",
				nulls: 0,
				err:   true,
			},
		},

		// TS fields

		{
			name:        "test invalid TS ( invalid string )",
			json:        createJson("tsField", "invalid TS"),
			index:       createIndex(false),
			len:         5,
			invalidJson: false,
			col: &Col{
				name:  "tsField",
				nulls: 0,
				err:   true,
			},
		},
		{
			name:        "test invalid TS (empty)",
			json:        createJson("tsField", ""),
			index:       createIndex(false),
			len:         5,
			invalidJson: false,
			col: &Col{
				name:  "tsField",
				nulls: 0,
				err:   true,
			},
		},
		{
			name:        "test invalid TS (nil)",
			json:        createJson("tsField", nil),
			index:       createIndex(false),
			len:         5,
			invalidJson: false,
			col: &Col{
				name:  "tsField",
				nulls: 0,
				err:   true,
			},
		},
		{
			name:        "test valid TS (string)",
			json:        createJson("tsField", strconv.Itoa(int(time.Now().UnixNano()))),
			index:       createIndex(false),
			len:         6,
			invalidJson: false,
			col: &Col{
				name:  "tsField",
				nulls: 0,
				err:   false,
			},
		},

		// nullable fields.

		{
			name:        "test nullable int",
			json:        createJson("intField", nil),
			index:       createIndex(true),
			len:         6,
			invalidJson: false,
			col: &Col{
				name:  "intField",
				nulls: 1,
				err:   false,
			},
		},
		{
			name:        "test nullable uint",
			json:        createJson("uintField", nil),
			index:       createIndex(true),
			len:         6,
			invalidJson: false,
			col: &Col{
				name:  "uintField",
				nulls: 1,
				err:   false,
			},
		}, {
			name:        "test nullable float",
			json:        createJson("floatField", nil),
			index:       createIndex(true),
			len:         6,
			invalidJson: false,
			col: &Col{
				name:  "floatField",
				nulls: 1,
				err:   false,
			},
		},
		{
			name:        "test nullable TS",
			json:        createJson("tsField", nil),
			index:       createIndex(true),
			len:         6,
			invalidJson: false,
			col: &Col{
				name:  "tsField",
				nulls: 1,
				err:   false,
			},
		},
		{
			name:        "test invalid json",
			json:        "this is an invalid json str",
			index:       createIndex(false),
			len:         0,
			invalidJson: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			in := NewIngester(test.index, bytes.NewBufferString(test.json))

			table, errors := in.IngestFromJSON()

			if test.invalidJson {
				assert.NotEmpty(t, errors, "a parse error should be returned")
				return
			}

			for _, f := range test.index.Fields {

				buf, found := table.Col(f.Id)

				if !assert.True(t, found, "column [%v] not found", f.Name) {
					return
				}

				if test.col != nil && test.col.name == f.Name {

					if f.Nullable {
						assert.Equal(t, test.col.nulls, buf.Nulls().GetCardinality())
					}

					if test.col.err {
						assert.True(t, containsCol(errors, test.col.name), "error msg expected")
					} else {
						assert.Falsef(t, containsCol(errors, test.col.name), "unexpected error msg %v", errors)
					}

				} else {
					if f.Nullable {
						assert.Equal(t, buf.Nulls().GetCardinality(), test.nulls)
					}
				}

				assert.Equal(t, test.len, buf.Len(), "wrong len value")
				assert.Equal(t, buf.Nullable(), f.Nullable)

			}

		})
	}
}

func containsCol(errors []IngestError, colName string) bool {
	for _, e := range errors {
		if e.Field == colName {
			return true
		}
	}
	return false
}

func createJson(field string, value interface{}) string {

	m := map[string]interface{}{
		"intField":    10,
		"floatField":  3.14159265359,
		"uintField":   20,
		"stringField": "stringfield",
		"textField":   "foo bar",
		"tsField":     time.Now().UnixNano(),
		"uuidField":   uuid.New().String(),
	}

	b, err := json.Marshal(m)

	if err != nil {
		panic(err)
	}

	if field != "" {
		if value == nil {
			delete(m, field)
		} else {
			m[field] = value
		}
	}

	b1, err := json.Marshal(m)

	if err != nil {
		panic(err)
	}

	s := string(b) + "\n"
	r := string(b1) + "\n"
	for i := 0; i < 5; i++ {
		r = r + s
	}

	return r

}
