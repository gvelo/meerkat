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

package tools

import (
	golangIo "io"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment/ondsk"
)

type DecoderIterator interface {
	HasNext() bool
	Err() error
	Next() interface{}
}

type Decoder interface {
	Decode(page *ondsk.ColumnIterator) DecoderIterator // iterates all values found.
}

// RLE
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

func NewRLEDecoder(buf []byte) *IntDecoder {
	return &IntDecoder{
		br: io.NewBinaryReader(buf),
	}
}

// Next returns true if a value was scanned.
func (d *IntDecoder) Next() interface{} {
	return d.Value
}

// Next returns true if a value was scanned.
func (d *IntDecoder) HasNext() bool {

	if d.Run > 1 {
		d.Run--
		return true
	}

	num, err := d.br.ReadVarUInt()
	if err == golangIo.EOF {
		return false
	}

	if err != nil {
		d.err = err
		return false
	}

	run, err := d.br.ReadVarUInt()
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

// DecodeInt encoded run.
func RLEIntDecode(buf []byte) (v []int, err error) {
	s := NewRLEDecoder(buf)

	for s.HasNext() {
		v = append(v, s.Next().(int))
	}

	return v, s.Err()
}

// Dictionary

type DictionaryDecoder struct {
}

func NewDictionaryDecoder() *DictionaryDecoder {
	return new(DictionaryDecoder)
}

type DictDecoder struct {
	br    *io.BinaryReader
	err   error
	dict  map[int]string
	Value string
}

// TODO: revisar porque seguro vamos a decodificar por pagina.
func (e *DictionaryDecoder) Decode(data interface{}) interface{} {
	/*
		pages := data.([]*ondsk.Page)

		p := pages[0]

		br := io.NewBinaryReader(p.Data)
		keys, _ := br.ReadVarUInt()
		dict := make(map[int]string)
		for i := 0; i < keys; i++ {
			k, _ := br.ReadVarUInt()
			v, _ := br.ReadString()
			dict[k] = v
		}

		d, _ := DictionaryDecode(p.Data[br.Offset:], dict)
		for i := 1; i < len(pages); i++ {
			dec, _ := DictionaryDecode(pages[i].Data, dict)
			d = append(d, dec...)
		} */

	return nil
}

// Err returns any error which ocurred during decoding.
func (d *DictDecoder) Err() error {
	return d.err
}

func NewDictDecoder(buf []byte, m map[int]string) *DictDecoder {
	return &DictDecoder{
		br:   io.NewBinaryReader(buf),
		dict: m,
	}
}

func (d *DictDecoder) Next() interface{} {
	return d.Value
}

// Next returns true if a value was scanned.
func (d *DictDecoder) HasNext() bool {

	num, err := d.br.ReadVarUInt()
	if err == golangIo.EOF {
		return false
	}

	if err != nil {
		d.err = err
		return false
	}

	d.Value = d.dict[num]
	return true
}

func DictionaryDecode(buf []byte, d map[int]string) (v []string, e error) {

	s := NewDictDecoder(buf, d)

	for s.HasNext() {
		v = append(v, s.Next().(string))
	}

	return v, s.Err()
}

// Snappy

type SnappyDecoder struct {
	Decoded string
	err     error
	value   string
	br      *io.BinaryReader
}

func NewSnappyDecoder() Decoder {
	return new(SnappyDecoder)
}

func (d *SnappyDecoder) Decode(page *ondsk.ColumnIterator) DecoderIterator {
	// s, _ := snappy.Decode(nil,page.value)
	return d
}

func (d *SnappyDecoder) Err() error {
	return nil
}

func (d *SnappyDecoder) HasNext() bool {

	return true
}

func (d *SnappyDecoder) Next() interface{} {
	return nil
}
