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
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/inmem"
	"meerkat/internal/storage/text"
	"meerkat/internal/tools"
	"path/filepath"
)

type HandlerFunc func(*MiddlewarePayload, ...interface{}) error

type Middleware func(HandlerFunc) HandlerFunc

func BuildChain(f HandlerFunc, m ...Middleware) HandlerFunc {
	if len(m) == 0 {
		return f
	}
	return m[0](BuildChain(f, m[1:cap(m)]...))
}

var RawEncoder = func(f HandlerFunc) HandlerFunc {
	return func(mp *MiddlewarePayload, args ...interface{}) error {
		fmt.Println("start Writer")
		r := make([]byte, 0)
		// nothing to do
		switch args[0].(type) {
		case []byte:
			r = args[0].([]byte)
		case []int:
			r = tools.UnsafeCastIntsToBytes(args[0].([]int))
		case []float64:
			r = tools.UnsafeCastFloatsToBytes(args[0].([]float64))
		case []string:
			r = tools.CastStringToBytes(args[0].([]string))
		}
		return f(mp, r)
	}
}

var BuildSkip = func(f HandlerFunc) HandlerFunc {
	return func(mp *MiddlewarePayload, args ...interface{}) error {

		mp.Posting = inmem.NewPostingStore()

		var comparator inmem.ComparatorInterface
		switch mp.Col.FieldInfo().Type {

		case segment.FieldTypeFloat:
			comparator = inmem.Float64Interface{}
		case segment.FieldTypeInt:
			comparator = inmem.IntInterface{}
		case segment.FieldTypeTimestamp:
			comparator = inmem.IntInterface{}
		}

		mp.Sl = inmem.NewSkipList(mp.Posting, comparator)

		for i := 0; i < mp.Col.Size(); i++ {
			x := mp.Col.Get(i)
			mp.Sl.Add(x, i)
		}
		mp.Cardinality = mp.Sl.Length
		return f(mp)
	}
}

var BuildBTrie = func(f HandlerFunc) HandlerFunc {
	return func(mp *MiddlewarePayload, args ...interface{}) error {

		col := mp.Col
		mp.Posting = inmem.NewPostingStore()
		mp.Trie = inmem.NewBtrie(mp.Posting)

		for i := 0; i < col.Size(); i++ {
			x := col.Get(i)
			mp.Trie.Add(x.(string), uint32(i)) // save val -> ids
		}
		mp.Cardinality = mp.Trie.Cardinality
		return f(mp)
	}
}

type MiddlewarePayload struct {
	Path        string
	Col         inmem.Column
	Trie        *inmem.BTrie
	Sl          *inmem.SkipList
	Posting     *inmem.PostingStore
	Pages       []*inmem.Page
	Cardinality int
}

func NewMiddlewarePayload(path string, column inmem.Column) *MiddlewarePayload {
	mp := new(MiddlewarePayload)
	mp.Path = path
	mp.Col = column
	return mp
}

func WriteToFile(mp *MiddlewarePayload, args ...interface{}) error {

	f := filepath.Join(mp.Path, mp.Col.FieldInfo().Name+idxPagExt)
	WriteStoreIdx(f, mp.Pages)

	if mp.Col.FieldInfo().Index {

		if mp.Col.FieldInfo().Type == segment.FieldTypeText {
			tokenizer := text.NewTokenizer()
			for i := 0; i < mp.Col.Size(); i++ {
				tokens := tokenizer.Tokenize(mp.Col.Get(i).(string))
				for _, token := range tokens {
					mp.Trie.Add(token, uint32(i))
				}
			}
		}

		f := filepath.Join(mp.Path, mp.Col.FieldInfo().Name+posExt)

		err := WritePosting(f, mp.Posting)
		if err != nil {
			return err
		}

		f = filepath.Join(mp.Path, mp.Col.FieldInfo().Name+idxPosExt)

		switch mp.Col.FieldInfo().Type {
		case segment.FieldTypeTimestamp, segment.FieldTypeFloat, segment.FieldTypeInt:
			err = WriteSkipList(f, mp.Sl)
			if err != nil {
				return err
			}
		case segment.FieldTypeKeyword, segment.FieldTypeText:
			err = WriteTrie(f, mp.Trie)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
