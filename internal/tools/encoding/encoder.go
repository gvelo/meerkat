package encoding

type Encoder interface {
	Encode(values interface{}) []byte
}

type RawEncoder struct {
	next Encoder
}

func (e *RawEncoder) Encode(values interface{}) []byte {
	r := make([]byte, 0)
	// nothing to do
	switch values.(type) {
	case []byte:
		r = values.([]byte)
	case []int:
		r = UnsafeCastIntsToBytes(values.([]int))
	case []float64:
		r = UnsafeCastFloatsToBytes(values.([]float64))
	case []string:
		r = CastStringToBytes(values.([]string))
	}
	if e.next != nil {
		return e.next.Encode(r)
	} else {
		return r
	}
}
