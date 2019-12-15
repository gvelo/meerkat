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
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"meerkat/internal/buffer"
	"meerkat/internal/schema"
	"time"
)

type IngestError struct {
	Line  int
	Field string
	Error string
}

type Ingester interface {
	IngestFromJSON() (*buffer.Table, []IngestError, error)
}

type ingester struct {
	log       zerolog.Logger
	indexInfo schema.IndexInfo
	reader    io.Reader
}

func NewIngester(index schema.IndexInfo, r io.Reader) *ingester {

	return &ingester{
		log:       log.With().Str("src", "ingester").Logger(),
		indexInfo: index,
		reader:    r,
	}

}

func (igt *ingester) IngestFromJSON() (*buffer.Table, []IngestError) {

	//TODO: CLOSE ?

	table := buffer.NewTable(igt.indexInfo)
	row := buffer.NewRow(len(igt.indexInfo.Fields))
	var errors []IngestError

	br := bufio.NewReader(igt.reader)

	// TODO(gvelo) configure MaxScanTokenSize
	scanner := bufio.NewScanner(br)

	line := 0

	for scanner.Scan() {

		line++

		b := scanner.Bytes()

		var i interface{}

		err := json.Unmarshal(b, &i)

		if err != nil {

			e := IngestError{
				Line:  line,
				Error: fmt.Sprintf("unable to parse json [%v]", err),
			}

			errors = append(errors, e)

			continue
		}

		jsonMap := i.(map[string]interface{})

		rowHasErrors := false

		for _, f := range igt.indexInfo.Fields {

			fv, ok := jsonMap[f.Name]

			if !ok {

				if f.Name == schema.TSFieldName {
					row.AddCol(schema.TSFieldName, int(time.Now().UnixNano()))
					continue
				}

				if f.Name == schema.IDFieldName {
					row.AddCol(schema.IDFieldName, uuid.New())
					continue
				}

				if !f.Nullable {

					err := fmt.Sprintf("field [%v] cannot be null", f.Name)

					ingestError := IngestError{
						Line:  line,
						Field: f.Name,
						Error: err,
					}

					errors = append(errors, ingestError)

					rowHasErrors = true
				}

				continue
			}

			v, err := convert(f.FieldType, fv)

			if err != nil {

				parseError := fmt.Errorf("cannot parse field [%v] %v", f.Name, err)

				ingestError := IngestError{
					Line:  line,
					Field: f.Name,
					Error: parseError.Error(),
				}

				errors = append(errors, ingestError)

				rowHasErrors = true

				continue
			}

			row.AddCol(f.Id, v)

		}

		if !rowHasErrors {
			table.AppendRow(row)
		}

		row.Reset()

	}

	if scanner.Err() != nil {

		line++

		e := IngestError{
			Line:  line,
			Error: scanner.Err().Error(),
		}

		errors = append(errors, e)

	}

	return table, errors

}
