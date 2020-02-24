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

package storage

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"meerkat/internal/buffer"
	"meerkat/internal/schema"
	"testing"
	"time"
)

func TestSegmentWriter_Write(t *testing.T) {

	index := schema.IndexInfo{
		Id:             "test-index",
		Name:           "test-index",
		Desc:           "test-index",
		Created:        time.Time{},
		Updated:        time.Time{},
		PartitionAlloc: schema.PartitionAlloc{},
		Fields: []schema.Field{
			{
				Id:        "_ts",
				Name:      "_ts",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_TIMESTAMP,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "_id",
				Name:      "_id",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_UUID,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "intFieldId",
				Name:      "intField",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_INT,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "stringFieldId",
				Name:      "stringField",
				Desc:      "",
				IndexId:   "test-index",
				FieldType: schema.FieldType_STRING,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
		},
	}

	table := buffer.NewTable(index)

	for i := 0; i < 10000; i++ {
		r := buffer.NewRow(4)
		r.AddCol("_ts", int(time.Now().UnixNano()))
		r.AddCol("_id", uuid.New())
		r.AddCol("intFieldId", rand.Int())
		r.AddCol("stringFieldId", fmt.Sprintf("row %v", i))
		table.AppendRow(r)
	}

	sw := NewSegmentWriter("/tmp/segment", uuid.New(), table)

	err := sw.Write()

	if err != nil {
		t.Error(err)
	}

}
