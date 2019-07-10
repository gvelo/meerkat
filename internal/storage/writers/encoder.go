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
	"log"
	"math"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/inmem"
	"path/filepath"
)

type Encoder interface {
	Encode(data interface{}) []*inmem.Page
}

var EncoderHandler = func(f HandlerFunc) HandlerFunc {
	return func(mp *MiddlewarePayload, args ...interface{}) error {

		var e Encoder
		switch mp.Col.FieldInfo().Type {
		case segment.FieldTypeFloat: // double delta (simple deberia... ) por ahora no encodea.
			e = NewFloatNOEncoder(mp)
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
		mp.Pages = e.Encode(mp.Col.Data())
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

func (e *SnappyEncoder) Encode(data interface{}) []*inmem.Page {

	src := data.([]string)

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
	for i := 0; i < len(src); i++ {
		// Encoding para despues separar los strings.
		str := fmt.Sprintf("%d %v", len(src[i]), src[i])
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
		pd.Total = len(src) - pd.StartID
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

func (e *RLEIntegerEncoder) Encode(data interface{}) []*inmem.Page {
	src := data.([]int)
	return EncodeRLE(e.path, e.info, src)
}

func EncodeRLE(path string, fieldInfo *segment.FieldInfo, data []int) []*inmem.Page {

	size := len(data)

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

	var cur = data[0]
	var run int

	for i := 0; i < size; i++ {
		num := data[i]

		if num != cur {

			pd.Enc = inmem.RLE
			pd.Offset = bw.Offset

			bw.WritePageHeader(pd)

			b := bw.Offset // to calculate payload
			bw.WriteVarInt(cur)
			bw.WriteVarInt(run)
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
	bw.WriteVarInt(cur)
	bw.WriteVarInt(run)
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

func (e *VarIntEncoder) Encode(data interface{}) []*inmem.Page {
	d := data.([]int)
	size := len(d)

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

	// create the batches
	batchSize := 1000
	var batches [][]int

	for batchSize < len(d) {
		d, batches = d[batchSize:], append(batches, d[0:batchSize:batchSize])
	}
	batches = append(batches, d)

	bt, _ := io.NewBufferBinaryWriter()

	idCounter := 0
	// write
	for i := 0; i < len(batches); i++ {

		pd := new(inmem.Page)
		pd.StartID = idCounter
		pd.Total = len(batches[i])
		idCounter = idCounter + pd.Total
		for x := 0; x < len(batches[i]); x++ {
			bt.WriteVarInt(batches[i][x])
		}
		pd.Offset = bw.Offset
		bw.WritePageHeader(pd)
		bt.Flush()
		pd.PayloadSize = bt.Offset
		bw.Write(bt.Buffer.Bytes())
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

func (e *EncodeDictionary) Encode(data interface{}) []*inmem.Page {

	src := data.([]string)
	size := len(src)

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
			if dict[src[i]] == 0 {
				dict[src[i]] = len(dict) + 1
			}
		}

		bt.WriteVarInt(dict[src[i]])

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
	bw.WriteVarInt(len(dict))
	for k, v := range dict {
		bw.WriteVarInt(v)
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

//TODO: implementar Double Delta (sin el double) o algo asi.
type FloatNoEncoder struct {
	path string
	info *segment.FieldInfo
}

func NewFloatNOEncoder(payload *MiddlewarePayload) Encoder {
	e := new(FloatNoEncoder)
	e.path = payload.Path
	e.info = payload.Col.FieldInfo()
	return e
}

func (e *FloatNoEncoder) Encode(data interface{}) []*inmem.Page {
	d := data.([]float64)
	size := len(d)

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

	// create the batches
	batchSize := 1000
	var batches [][]float64

	for batchSize < len(d) {
		d, batches = d[batchSize:], append(batches, d[0:batchSize:batchSize])
	}
	batches = append(batches, d)

	bt, _ := io.NewBufferBinaryWriter()

	// write
	for i := 0; i < len(batches); i++ {

		pd := new(inmem.Page)
		pd.StartID = i
		pd.Total = len(batches[i])
		for x := 0; x < len(batches[i]); x++ {
			x := math.Float64bits(batches[i][x])
			bt.WriteFixedUint64(x)
		}
		pd.Offset = bw.Offset
		bw.WritePageHeader(pd)
		bt.Flush()
		pd.PayloadSize = bt.Offset
		bw.Write(bt.Buffer.Bytes())
		pages = append(pages, pd)
	}

	return pages
}
