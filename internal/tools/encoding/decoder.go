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

import (
	"github.com/golang/snappy"
	golangIo "io"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/inmem"
	"meerkat/internal/tools/utils"
)

type Decoder interface {
	Decode(data interface{}) interface{}
}

// TODO Capas convenga hacer un

var DecoderHandler = func(f HandlerFunc) HandlerFunc {
	return func(args ...interface{}) {
		col := args[0].(inmem.Column)
		slice := args[1].(utils.Slicer)
		//total := slice.Total()
		//c := args[2].(int)

		var e Encoder
		var pages []*inmem.PageDescriptor

		switch col.FieldInfo().Type {
		//case segment.FieldTypeFloat:
		//	e = inmem.Float64Interface{}
		//case segment.FieldTypeKeyword:
		//	e = inmem.Float64Interface{}
		case segment.FieldTypeText:
			e = NewSnappyEncoder()
			pages = e.Encode(slice)
		case segment.FieldTypeInt:
			e = NewRLEEncoder()
			pages = e.Encode(slice)
		//	e = inmem.IntInterface{}
		//case segment.FieldTypeTimestamp:
		//	e = inmem.IntInterface{}
		default:
			print("Could not encode.")
			return

		}

		f(pages)

	}

}

type SnappyDecoder struct {
}

func NewSnappyDecoder() Decoder {
	return new(SnappyDecoder)
}

func (d *SnappyDecoder) Decode(data interface{}) interface{} {
	e := data.([]byte)
	s, _ := snappy.Decode(nil, e)
	return string(s)
}

type RLEIntegerDecoder struct {
}

func NewRLEDecoder() Decoder {
	return new(RLEIntegerDecoder)
}

func (e *RLEIntegerDecoder) Decode(data interface{}) interface{} {
	src := data.([]byte)
	d, _ := Decode(src)
	return d
}

type IntDecoder struct {
	Value int
	Run   int
	br    *io.BinaryReader
	err   error
}

// Err returns any error which ocurred during decoding.
func (d *IntDecoder) Err() error {
	return d.err
}

func NewIntDecoder(buf []byte) *IntDecoder {
	return &IntDecoder{
		br: io.NewBinaryReader(buf),
	}
}

// Next returns true if a value was scanned.
func (d *IntDecoder) Next() bool {

	if d.Run > 1 {
		d.Run--
		return true
	}

	num, err := d.br.ReadVarInt()
	if err == golangIo.EOF {
		return false
	}

	if err != nil {
		d.err = err
		return false
	}

	run, err := d.br.ReadVarInt()
	if err == golangIo.EOF {
		d.err = golangIo.ErrUnexpectedEOF
		return false
	}

	if err != nil {
		d.err = err
		return false
	}

	d.Value = num
	d.Run = run

	return true
}

// DecodeInt64 encoded run.
func Decode(buf []byte) (v []int, err error) {
	s := NewIntDecoder(buf)

	for s.Next() {
		v = append(v, s.Value)
	}

	return v, s.Err()
}
