package writers

import (
	"eventdb/readers"
	"eventdb/segment"
	"testing"
)

func TestReadWritePosting(t *testing.T) {

	posting := make([]*segment.PostingList, 0)

	for i := 0; i < 1000; i++ {
		p := segment.NewPostingList(uint32(i))
		posting = append(posting, p)
	}

	file := "/tmp/posting.bin"
	err := WritePosting(file, posting)

	if err != nil {
		t.Error(err)
	}

	pr, err := readers.NewPostingReader(file)

	if err != nil {
		t.Error(err)
		return
	}

	for i, p := range posting {

		b, err := pr.Read(p.Offset)

		if err != nil {
			t.Error(err)
		}

		if !b.ContainsInt(i) || b.GetCardinality() != 1 {
			t.Errorf("Bitmap doesn't contain expected values.")
		}

	}

}
