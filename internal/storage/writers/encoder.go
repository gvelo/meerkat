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

package writers

import (
	"fmt"
	"github.com/golang/snappy"
	"github.com/spenczar/fpc"
	"log"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/inmem"
	"path/filepath"
)

type Encoder interface {
	Encode(col inmem.Column) []*inmem.Page
}

var EncoderHandler = func(f HandlerFunc) HandlerFunc {
	return func(mp *MiddlewarePayload, args ...interface{}) error {

		var e Encoder
		switch mp.Col.FieldInfo().Type {
		case segment.FieldTypeFloat: // double delta (simple deberia... ) por ahora no encodea.
			e = NewFloatFPCEncoder(mp)
		case segment.FieldTypeKeyword: // dict
			e = NewDictionaryEncoder(mp)
		case segment.FieldTypeText: // Snappy
			e = NewSnappyEncoder(mp)
		case segment.FieldTypeInt: //  RLE, Simple8b, Varint.
			e = selectIntEncoder(mp)
		case segment.FieldTypeTimestamp: // RLE, Simple8b, Varint.
			e = selectIntEncoder(mp)
		default:
			print("Could not encode.")
			return nil

		}
		mp.Pages = e.Encode(mp.Col)
		return f(mp)

	}

}

func selectIntEncoder(mp *MiddlewarePayload) Encoder {
	c := mp.Sl.Length
	var e Encoder
	if c < 10000 {
		e = NewRLEEncoder(mp)
	} else {
		e = NewVarIntEncoder(mp)
	}
	return e
}

type SnappyEncoder struct {
	path      string
	fieldInfo *segment.FieldInfo
}

func NewSnappyEncoder(mp *MiddlewarePayload) Encoder {
	e := new(SnappyEncoder)
	e.path = mp.Path
	e.fieldInfo = mp.Col.FieldInfo()
	return e
}

func (e *SnappyEncoder) Encode(col inmem.Column) []*inmem.Page {

	f := filepath.Join(e.path, e.fieldInfo.Name+pagExt)

	bw, _ := io.NewBinaryWriter(f)
	defer bw.Close()

	err := bw.WriteHeader(io.RowStoreV1)
	if err != nil {
		return nil
	}

	pages := make([]*inmem.Page, 0)
	b := make([]byte, 0, 40000)
	pd := new(inmem.Page)
	pd.StartID = 0

	for i := 0; i < col.Size(); i++ {
		// Encoding para despues separar los strings.
		str := fmt.Sprintf("%d %v", len(col.Get(i).(string)), col.Get(i))
		b = append(b, []byte(str)...)

		if len(b) >= 32000 {
			ed := snappy.Encode(nil, b)
			pd.Enc = inmem.Snappy
			pd.PayloadSize = len(ed)
			log.Printf("pagina %v len(b) %v", pd, len(b))

			pd.Offset = bw.Offset
			pd.Total = i - pd.StartID
			bw.WritePageHeader(pd)
			bw.Write(ed)

			pages = append(pages, pd)
			pd = new(inmem.Page)

			pd.StartID = i
			b = make([]byte, 0, 40000)

		}
	}

	if len(b) > 0 {
		ed := snappy.Encode(nil, b)
		pd.Enc = inmem.Snappy
		pd.PayloadSize = len(ed)
		log.Printf("pagina %v len(b) %v", pd, len(b))

		pd.Offset = bw.Offset
		pd.Total = col.Size() - pd.StartID
		bw.WritePageHeader(pd)
		bw.Write(ed)

		pages = append(pages, pd)
	}

	return pages
}

type RLEIntegerEncoder struct {
	path string
	info *segment.FieldInfo
}

func NewRLEEncoder(payload *MiddlewarePayload) Encoder {
	e := new(RLEIntegerEncoder)
	e.path = payload.Path
	e.info = payload.Col.FieldInfo()
	return e
}

func (e *RLEIntegerEncoder) Encode(col inmem.Column) []*inmem.Page {
	return EncodeRLE(e.path, e.info, col)
}

func EncodeRLE(path string, fieldInfo *segment.FieldInfo, col inmem.Column) []*inmem.Page {

	size := col.Size()

	if size == 0 {
		return nil
	}

	pages := make([]*inmem.Page, 0)

	pd := new(inmem.Page)
	pd.StartID = 0

	f := filepath.Join(path, fieldInfo.Name+pagExt)

	bw, _ := io.NewBinaryWriter(f)
	defer bw.Close()

	err := bw.WriteHeader(io.RowStoreV1)
	if err != nil {
		return nil
	}

	var cur = col.Get(0).(int)
	var run int

	for i := 0; i < size; i++ {
		num := col.Get(i).(int)

		if num != cur {

			pd.Enc = inmem.RLE
			pd.Offset = bw.Offset

			bw.WritePageHeader(pd)

			b := bw.Offset // to calculate payload
			bw.WriteVarUInt(cur)
			bw.WriteVarUInt(run)
			a := bw.Offset
			pd.PayloadSize = a - b
			pd.Total = i - pd.StartID
			pages = append(pages, pd)

			cur = num
			run = 0

			pd = new(inmem.Page)
			pd.StartID = i

		}

		run++
	}

	pd.Enc = inmem.RLE
	pd.Offset = bw.Offset
	pd.Total = size - pd.StartID
	bw.WritePageHeader(pd)
	bw.WriteVarUInt(cur)
	bw.WriteVarUInt(run)
	pages = append(pages, pd)

	return pages
}

type VarIntEncoder struct {
	path string
	info *segment.FieldInfo
}

func NewVarIntEncoder(payload *MiddlewarePayload) Encoder {
	e := new(VarIntEncoder)
	e.path = payload.Path
	e.info = payload.Col.FieldInfo()
	return e
}

func (e *VarIntEncoder) Encode(col inmem.Column) []*inmem.Page {
	size := col.Size()

	if size == 0 {
		return nil
	}

	pages := make([]*inmem.Page, 0)

	f := filepath.Join(e.path, e.info.Name+pagExt)

	bw, _ := io.NewBinaryWriter(f)
	defer bw.Close()

	err := bw.WriteHeader(io.RowStoreV1)
	if err != nil {
		return nil
	}

	bt, _ := io.NewBufferBinaryWriter()

	bz := 1000
	idCounter := 0

	// write
	for i := 0; i < col.Size(); i++ {

		bt.WriteVarUInt(col.Get(i).(int))

		if i > 0 && i%bz == 0 {
			pd := new(inmem.Page)
			pd.Total = bz
			pd.StartID = idCounter
			idCounter = idCounter + bz
			pd.Offset = bw.Offset
			pd.PayloadSize = bt.Offset
			bw.WritePageHeader(pd)
			bt.Flush()
			bw.Write(bt.Buffer.Bytes())
			bt.Buffer.Reset()
			pages = append(pages, pd)
		}

	}

	if bt.Offset > 0 {
		pd := new(inmem.Page)
		pd.Total = bz
		pd.StartID = idCounter
		idCounter = idCounter + bz
		pd.Offset = bw.Offset
		pd.PayloadSize = bt.Offset
		bw.WritePageHeader(pd)
		bt.Flush()
		bw.Write(bt.Buffer.Bytes())
		bt.Buffer.Reset()
		pages = append(pages, pd)
	}

	return pages
}

type EncodeDictionary struct {
	path        string
	info        *segment.FieldInfo
	cardinality int
}

func NewDictionaryEncoder(mp *MiddlewarePayload) Encoder {
	e := new(EncodeDictionary)
	e.path = mp.Path
	e.info = mp.Col.FieldInfo()
	e.cardinality = mp.Cardinality
	return e
}

func (e *EncodeDictionary) Encode(col inmem.Column) []*inmem.Page {

	size := col.Size()

	if size == 0 {
		return nil
	}

	f := filepath.Join(e.path, e.info.Name+pagExt)

	bw, _ := io.NewBinaryWriter(f)
	defer bw.Close()

	pages := make([]*inmem.Page, 0)

	pd := new(inmem.Page)
	pd.Enc = inmem.Dictionary
	pd.StartID = 0

	bt, _ := io.NewBufferBinaryWriter()
	dict := make(map[string]int)

	for i := 0; i < size; i++ {
		if len(dict) < e.cardinality {
			if dict[col.Get(i).(string)] == 0 {
				dict[col.Get(i).(string)] = len(dict) + 1
			}
		}

		bt.WriteVarUInt(dict[col.Get(i).(string)])

		if i != 0 && i%1000 == 0 {

			pd.Offset = bw.Offset
			pd.Total = i - pd.StartID
			bw.WritePageHeader(pd)
			bw.Write(buildData(bt))
			pages = append(pages, pd)
			pd = new(inmem.Page)
			pd.Enc = inmem.Dictionary
			pd.StartID = i
		}

	}

	pd.Offset = bw.Offset
	pd.Total = size - pd.StartID
	bw.WritePageHeader(pd)
	bw.Write(buildData(bt))

	pages = append(pages, pd)

	// write dictionary to file´s bottom
	offset := bw.Offset
	bw.WriteVarUInt(len(dict))
	for k, v := range dict {
		bw.WriteVarUInt(v)
		bw.WriteBytes([]byte(k))
	}
	bw.WriteFixedUint64(uint64(offset))

	return pages
}

func buildData(bw *io.BufferBinaryWriter) []byte {
	r := make([]byte, 0, bw.Offset)
	bw.Flush() // Writes to buffer.
	r = append(r, bw.Buffer.Bytes()...)
	bw.Buffer.Reset()
	return r

}

type FPCEncoder struct {
	path string
	info *segment.FieldInfo
}

func NewFloatFPCEncoder(payload *MiddlewarePayload) Encoder {
	e := new(FPCEncoder)
	e.path = payload.Path
	e.info = payload.Col.FieldInfo()
	return e
}

func (e *FPCEncoder) Encode(col inmem.Column) []*inmem.Page {
	size := col.Size()

	if size == 0 {
		return nil
	}

	pages := make([]*inmem.Page, 0)

	f := filepath.Join(e.path, e.info.Name+pagExt)

	bw, _ := io.NewBinaryWriter(f)
	defer bw.Close()

	err := bw.WriteHeader(io.RowStoreV1)
	if err != nil {
		return nil
	}

	bsz := 1000
	idCounter := 0
	t := 0

	bt, _ := io.NewBufferBinaryWriter()
	fpcw := fpc.NewWriter(bt)
	// write
	for i := 0; i < col.Size(); i++ {

		fpcw.WriteFloat(col.Get(i).(float64))
		t = t + 1

		if i > 0 && i%bsz == 0 {

			pd := new(inmem.Page)
			pd.Total = t
			pd.StartID = idCounter
			idCounter = idCounter + pd.Total
			pd.Offset = bw.Offset
			pd.Enc = inmem.FPC
			fpcw.Flush()
			bt.Flush()
			pd.PayloadSize = bt.Offset
			bw.WritePageHeader(pd)
			bw.Write(bt.Buffer.Bytes())
			pages = append(pages, pd)
			t = 0

		}

	}

	if t > 0 {
		pd := new(inmem.Page)
		pd.Total = t
		pd.StartID = idCounter
		idCounter = idCounter + pd.Total
		pd.Offset = bw.Offset
		pd.Enc = inmem.FPC
		fpcw.Flush()
		bt.Flush()
		pd.PayloadSize = bt.Offset
		bw.WritePageHeader(pd)
		bw.Write(bt.Buffer.Bytes())
		pages = append(pages, pd)
	}

	return pages
}
