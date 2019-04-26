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

	defer bw.Close()

	err = bw.WriteHeader(io.PostingListV1)

	if err != nil {
		return err
	}

	for _, p := range posting {
			p.Bitmap.RunOptimize()
		p.Offset = bw.Offset
		_, err := p.Bitmap.WriteTo(bw)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil

}
