package writers

import (
	"eventdb/io"
	"eventdb/segment/inmem"
)

func WritePosting(name string, posting []*inmem.PostingList) error {

	bw, err := io.NewBinaryWriter(name)

	if err != nil {
		return err
	}

	bw.WriteHeader(io.PostingListV1)

	for _, p := range posting {
		n, err := p.Bitmap.WriteTo(bw)
	}

}
