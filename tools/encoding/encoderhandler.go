package encoding

import (
	"eventdb/segment"
	"eventdb/segment/inmem"
)

type EncoderHandler struct {
	root Encoder
}

func NewEncoderHandler(fieldInfo *segment.FieldInfo, page *inmem.Page) *EncoderHandler {
	// create a chain of encoders
	chain := &RawEncoder{next: nil}
	return &EncoderHandler{root: chain}
}

func (e *EncoderHandler) DoEncode(slice interface{}) interface{} {
	return e.root.Encode(slice)
}
