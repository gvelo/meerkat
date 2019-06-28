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
)

type HandlerFunc func(...interface{})

type middleware func(HandlerFunc) HandlerFunc

func BuildChain(f HandlerFunc, m ...middleware) HandlerFunc {
	if len(m) == 0 {
		return f
	}
	return m[0](BuildChain(f, m[1:cap(m)]...))
}

func main() {

	var privateChain = []middleware{
		Encoding,
		Writer,
	}

	chain := BuildChain(WriteToFile, privateChain...)
	chain("3", "3")

}

var Encoding = func(f HandlerFunc) HandlerFunc {
	return func(args ...interface{}) {
		fmt.Println("start Encoding")
		f(args)
		fmt.Println("end Encoding")
	}
}

var Writer = func(f HandlerFunc) HandlerFunc {
	return func(args ...interface{}) {
		fmt.Println("start Writer")
		f(args)
		fmt.Println("end Writer")
	}
}

var RawEncoder = func(f HandlerFunc) HandlerFunc {
	return func(args ...interface{}) {
		fmt.Println("start Writer")
		r := make([]byte, 0)
		// nothing to do
		switch args[0].(type) {
		case []byte:
			r = args[0].([]byte)
		case []int:
			r = UnsafeCastIntsToBytes(args[0].([]int))
		case []float64:
			r = UnsafeCastFloatsToBytes(args[0].([]float64))
		case []string:
			r = CastStringToBytes(args[0].([]string))
		}
		f(r)
	}
}

// this is the handler func we are wrapping with middlewares
func WriteToFile(args ...interface{}) {

	println("w,r")
}
