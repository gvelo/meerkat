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
	"fmt"
	"github.com/golang/snappy"
	"log"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/inmem"
	"meerkat/internal/tools/utils"
)

type Encoder interface {
	Encode(data interface{}) []*inmem.PageDescriptor
}

var EncoderHandler = func(f HandlerFunc) HandlerFunc {
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

type SnappyEncoder struct {
}

func NewSnappyEncoder() Encoder {
	return new(SnappyEncoder)
}

func (e *SnappyEncoder) Encode(data interface{}) []*inmem.PageDescriptor {

	src := data.(utils.Slicer).Get().([]string)

	pages := make([]*inmem.PageDescriptor, 0)
	b := make([]byte, 0, 40000)
	pd := new(inmem.PageDescriptor)
	pd.StartID = 0
	for i := 0; i < len(src); i++ {
		// Encoding para despues separar los strings.
		str := fmt.Sprintf("%d %v", len(src[i]), src[i])
		b = append(b, []byte(str)...)

		if len(b) >= 32000 || i == len(src) {
			pd.Data = snappy.Encode(nil, b)
			pd.Enc = inmem.Snappy

			log.Printf("pagina %v len(b) %v", pd, len(b))

			pages = append(pages, pd)
			pd = new(inmem.PageDescriptor)
			pd.Data = nil
			pd.StartID = i
			b = make([]byte, 0, 40000)

		}
	}
	return pages
}

type RLEIntegerEncoder struct {
}

func NewRLEEncoder() Encoder {
	return new(RLEIntegerEncoder)
}

func (e *RLEIntegerEncoder) Encode(data interface{}) []*inmem.PageDescriptor {

	src := data.(utils.Slicer).Get().([]int)
	return EncodeRLE(src)
}

func EncodeRLE(nums []int) []*inmem.PageDescriptor {
	size := len(nums)

	if size == 0 {
		return nil
	}

	pages := make([]*inmem.PageDescriptor, 0)

	pd := new(inmem.PageDescriptor)
	pd.StartID = 0

	bw, _ := io.NewBufferBinaryWriter()
	var cur = nums[0]
	var run int

	for i := 0; i < size; i++ {
		num := nums[i]

		if num != cur {
			bw.WriteVarInt(cur)
			bw.WriteVarInt(run)
			cur = num
			run = 0

			bw.Flush()
			pd.Data = bw.Buffer.Bytes()
			pd.Enc = inmem.RLE
			pages = append(pages, pd)

			pd = new(inmem.PageDescriptor)
			pd.StartID = i

		}

		run++
	}

	bw.WriteVarInt(cur)
	bw.WriteVarInt(run)

	bw.Flush() // Writes to buffer.
	pd.Data = bw.Buffer.Bytes()
	pd.Enc = inmem.RLE
	pages = append(pages, pd)
	return pages
}
