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

package testutil

import (
	"math/rand"
	"strings"
)

var chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
	"abcdefghijklmnopqrstuvwxyzåäö" +
	"0123456789")

func RandomString(maxLenght int) string {
	length := rand.Intn(maxLenght) + 1
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func RandomBytes(l int) []byte {
	i := rand.Intn(l) + 1
	b := make([]byte, i)
	rand.Read(b)
	return b
}
