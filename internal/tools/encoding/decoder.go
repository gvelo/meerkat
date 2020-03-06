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

package encoding

type Decoder interface {
	Decode(values interface{}) []byte
}

type RawDecoder struct {
	next Decoder
}

func (e *RawDecoder) Decoder(values []byte, t interface{}) interface{} {
	var r interface{}
	// nothing to do
	switch t.(type) {
	case []byte:
		r = values
	case []int:
		r = UnsafeCastByteSliceToIntSlice(values)
	case []float64:
		r = UnsafeCastByteSliceToFloatSlice(values)
	case []string:
		r = CastByteSliceToStringSlice(values)
	}
	if e.next != nil {
		return e.next.Decode(r)
	} else {
		return r
	}
}
