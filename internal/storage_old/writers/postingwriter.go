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
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment/inmem"
)

func WritePosting(name string, posting *inmem.PostingStore) error {

	bw, err := io.NewBinaryWriter(name)

	if err != nil {
		return err
	}

	defer bw.Close()

	err = bw.WriteHeader(io.PostingListV1)

	if err != nil {
		return err
	}

	for _, p := range posting.Store {
		p.Bitmap.RunOptimize()
		p.Offset = bw.Offset
		_, err := p.Bitmap.WriteTo(bw)
		if err != nil {
			return err
		}
	}

	return nil

}
