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

// +build darwin linux

package mmap

import (
	"os"
	"syscall"
)

func mmap(f *os.File, size int64) ([]byte, error) {

	b, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)

	// TODO add kernel advicing.

	if err != nil {
		return nil, err
	}

	return b, nil

}

func unmap(ref []byte) error {
	return syscall.Munmap(ref)
}
