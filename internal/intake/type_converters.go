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
	"fmt"
	"github.com/google/uuid"
	"meerkat/internal/schema"
	"strconv"
	"time"
)

// TODO(gvelo): remove invalid int values, json parser only emit string and floats.

func convert(toType schema.FieldType, value interface{}) (interface{}, error) {
	switch toType {
	case schema.FieldType_FLOAT:
		return ToFloat(value)
	case schema.FieldType_INT:
		return ToInt(value)
	case schema.FieldType_UINT:
		return ToUint(value)
	case schema.FieldType_STRING:
		return ToString(value)
	case schema.FieldType_TEXT:
		return ToString(value)
	case schema.FieldType_TIMESTAMP:
		return ToDate(value)
	case schema.FieldType_UUID:
		return ToUUID(value)

	default:
		panic(fmt.Sprintf("invalid field type [%v]", toType))
	}
}

func ToFloat(value interface{}) (interface{}, error) {
	switch value.(type) {
	case string:
		return strconv.ParseFloat(value.(string), 64)
	case float64:
		return value, nil
	case int:
		return float64(value.(int)), nil
	case uint:
		return float64(value.(uint)), nil
	default:
		err := fmt.Errorf("cannot convert from %T to Float", value)
		return nil, err
	}
}

func ToInt(value interface{}) (interface{}, error) {
	switch value.(type) {
	case string:
		v, err := strconv.ParseInt(value.(string), 0, 64)
		if err != nil {
			return nil, err
		}
		return int(v), nil
	case float64:
		return int(value.(float64)), nil
	case int:
		return value, nil
	case uint:
		return int(value.(uint)), nil
	default:
		err := fmt.Errorf("cannot convert from %T to int", value)
		return nil, err

	}
}

func ToUint(value interface{}) (interface{}, error) {
	switch value.(type) {
	case string:
		v, err := strconv.ParseUint(value.(string), 0, 64)
		if err != nil {
			return nil, err
		}
		return uint(v), nil
	case float64:
		v := value.(float64)
		if v < 0 {
			return nil, fmt.Errorf("cannot convert [%v] to uint", v)
		}
		return uint(v), nil
	case int:
		v := value.(int)
		if v < 0 {
			return nil, fmt.Errorf("cannot convert [%v] to uint", v)
		}
		return v, nil
	case uint:
		return value, nil
	default:
		err := fmt.Errorf("cannot convert from %T to uint", value)
		return nil, err

	}
}

func ToString(value interface{}) (interface{}, error) {
	switch value.(type) {
	case string:
		return value, nil
	case float64:
		return strconv.FormatFloat(value.(float64), 'E', -1, 64), nil
	case int:
		return strconv.Itoa(value.(int)), nil
	case uint:
		return strconv.FormatUint(value.(uint64), 10), nil
	default:
		err := fmt.Errorf("cannot convert from %T to string", value)
		return nil, err

	}
}

func ToUUID(value interface{}) (interface{}, error) {
	switch value.(type) {
	case string:
		uid, err := uuid.Parse(value.(string))
		return uid, err
	default:
		err := fmt.Errorf("cannot convert from %T to UUID", value)
		return nil, err

	}
}

func ToDate(value interface{}) (interface{}, error) {
	switch value.(type) {
	case string:
		n, err := strconv.Atoi(value.(string))
		if err == nil {
			return n, nil
		}
		t, err := time.Parse(time.RFC3339Nano, value.(string))
		if err != nil {
			return nil, err
		}
		return int(t.UnixNano()), nil
	case float64:
		return int(value.(float64)), nil
	default:
		err := fmt.Errorf("cannot convert from %T to Date", value)
		return nil, err
	}
}
