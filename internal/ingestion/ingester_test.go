// Copyright 2020 The Meerkat Authors
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

package ingestion

import (
	"bytes"
	"fmt"
	"testing"
)

var s string

func init() {
	s := `{"_ts":"2020-05-11T18:46:06.577Z","columnA":23}
          {"_ts":"2020-05-11T18:46:06.577Z","columnA":23}
          {"_ts":"2020-05-11T18:46:06.577Z","columnA":23}
`
	for i := 0; i < 20; i++ {
		s=s+s
	}

	fmt.Println(len(s))
}

func TestTest(t *testing.T) {


//fmt.Println(s)
	r := bytes.NewReader([]byte(s))

	ing := NewIngester()
	_,err := ing.Parse(r,"testTable",3)

	fmt.Println(err)
	//fmt.Println(table.partitions[1].colSize[2])
	//fmt.Println(string(table.partitions[1].writer.Buf.Data()))
	//
	//for _, d := range table.partitions[1].writer.Buf.Data() {
	//	fmt.Println(d,string(d))
	//}


}


