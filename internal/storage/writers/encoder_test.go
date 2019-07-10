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
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/inmem"
	"meerkat/internal/tools/utils"
	"testing"
	"time"
)

func TestSnappyEncoder_Encode(t *testing.T) {

	assert := assert.New(t)

	fi := &segment.FieldInfo{0, "_time", segment.FieldTypeTimestamp, true}
	c := inmem.NewColumnt(fi)
	mp := NewMiddlewarePayload("/tmp/", c, &utils.SlicerString{})
	enc := NewSnappyEncoder(mp)

	slice := make([]string, 0)

	for i := 0; i < 10000; i++ {
		slice = append(slice, fmt.Sprintf("String numero %d", i))
	}

	s := enc.Encode(slice)

	sum := 0
	for i := 0; i < len(s); i++ {
		sum = sum + s[i].Total
		assert.NotNil(s[i].PayloadSize)
		assert.NotNil(inmem.Snappy, s[i].Enc)

	}

	assert.Equal(10000, sum)

}

func TestRLEIntegerEncoder_Encode(t *testing.T) {

	assert := assert.New(t)

	fi := &segment.FieldInfo{0, "Rle", segment.FieldTypeInt, true}
	c := inmem.NewColumnt(fi)

	slice := make([]int, 0)
	x := 0
	for i := 0; i < 10000; i++ {
		if i%1000 == 0 {
			x = x + 1
		}
		slice = append(slice, x)
	}

	mp := NewMiddlewarePayload("/tmp/", c, &utils.SlicerInt{})
	enc := NewRLEEncoder(mp)
	s := enc.Encode(slice)

	assert.NotNil(s)
	assert.Equal(10, len(s))

	sum := 0
	for i := 0; i < len(s); i++ {
		sum = sum + s[i].Total
		assert.NotNil(s[i].PayloadSize)
		assert.NotNil(inmem.RLE, s[i].Enc)
	}

	assert.Equal(10000, sum)

}

func TestNewDictionaryEncoder(t *testing.T) {

	assert := assert.New(t)

	fi := &segment.FieldInfo{0, "dict", segment.FieldTypeKeyword, true}
	c := inmem.NewColumnt(fi)

	dict := []string{"hola", "ma", "no", "la"}
	x := 0
	s := dict[x]
	for i := 0; i < 10000; i++ {
		s = dict[x]
		if i > 0 && i%(10000/4) == 0 {
			x++
		}
		c.Add(s)
	}

	mp := NewMiddlewarePayload("/tmp/", c, new(utils.SlicerInt))
	mp.Cardinality = len(dict)
	e := NewDictionaryEncoder(mp)
	pages := e.Encode(c.Data())

	assert.NotNil(pages)
	assert.Equal(10, len(pages))

	sum := 0
	for i := 0; i < len(pages); i++ {
		sum = sum + pages[i].Total
		assert.NotNil(pages[i].PayloadSize)
		assert.NotNil(inmem.Dictionary, pages[i].Enc)
	}

	assert.Equal(10000, sum)

}

func TestColumnTimeStamp_Middleware(t *testing.T) {

	a := assert.New(t)

	fi := &segment.FieldInfo{0, "_time", segment.FieldTypeTimestamp, true}
	c := inmem.NewColumnt(fi)

	for i := 0; i < 10000; i++ {
		c.Add(int(time.Now().Nanosecond()) + i)
	}

	var privateChain = []Middleware{
		BuildSkip,
		EncoderHandler,
	}

	mp := NewMiddlewarePayload("/tmp/", c, new(utils.SlicerInt))

	executeChain := BuildChain(WriteToFile, privateChain...)
	err := executeChain(mp)
	a.Equal(nil, err)

	sum := 0
	for i := 0; i < len(mp.Pages); i++ {
		sum = sum + mp.Pages[i].Total
		a.NotNil(mp.Pages[i].PayloadSize)
		a.Equal(inmem.RLE, mp.Pages[i].Enc)
	}

	a.Equal(10000, sum)

}

func TestColumnInt_Middleware(t *testing.T) {

	a := assert.New(t)

	fi := &segment.FieldInfo{0, "num", segment.FieldTypeInt, true}
	c := inmem.NewColumnt(fi)

	for i := 0; i < 10000; i++ {
		c.Add(i)
	}

	var privateChain = []Middleware{
		BuildSkip,
		EncoderHandler,
	}

	mp := NewMiddlewarePayload("/tmp/", c, new(utils.SlicerInt))

	executeChain := BuildChain(WriteToFile, privateChain...)
	err := executeChain(mp)
	a.Equal(nil, err)

	sum := 0
	for i := 0; i < len(mp.Pages); i++ {
		sum = sum + mp.Pages[i].Total
		a.NotNil(mp.Pages[i].PayloadSize)
	}

	a.Equal(10000, sum)

}

func TestColumnText_Middleware(t *testing.T) {

	a := assert.New(t)

	fi := &segment.FieldInfo{0, "text", segment.FieldTypeText, true}
	c := inmem.NewColumnt(fi)

	for i := 0; i < 100000; i++ {
		c.Add(fmt.Sprintf("Texto %d", i))
	}

	var privateChain = []Middleware{
		BuildBTrie,
		EncoderHandler,
	}

	mp := NewMiddlewarePayload("/tmp/", c, new(utils.SlicerString))

	executeChain := BuildChain(WriteToFile, privateChain...)
	err := executeChain(mp)
	a.Equal(nil, err)

	sum := 0
	for i := 0; i < len(mp.Pages); i++ {
		sum = sum + mp.Pages[i].Total
	}

	a.Equal(100000, sum)

}
