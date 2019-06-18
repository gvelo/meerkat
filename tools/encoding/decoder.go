package encoding

type Decoder interface {
	Decode(values interface{}) []byte
}

type RawDecoder struct {
	next Decoder
}

func (e *RawDecoder) Decoder(values []byte, t interface{}) interface{} {
	var r interface{}
	// nothing to do
	switch t.(type) {
	case []byte:
		r = values
	case []int:
		r = UnsafeCastBytesToInts(values)
	case []float64:
		r = UnsafeCastBytesToFloats(values)
	case []string:
		r = CastBytesToString(values)
	}
	if e.next != nil {
		return e.next.Decode(r)
	} else {
		return r
	}
}
