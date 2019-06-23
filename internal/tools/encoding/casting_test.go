package encoding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnsafeCastBytesToInt64(t *testing.T) {

	assert := assert.New(t)

	list := make([]int, 0)

	for i := 0; i < 1000; i++ {
		list = append(list, i)
	}

	before := make([]int, 1000)
	copy(before, list)

	byteList := UnsafeCastIntsToBytes(list)
	res := UnsafeCastBytesToInts(byteList)

	assert.Equal(before, res)

}

func TestUnsafeCastBytesToFloat(t *testing.T) {

	assert := assert.New(t)

	list := make([]float64, 0)

	for i := 0; i < 1000; i++ {
		list = append(list, float64(i))
	}

	before := make([]float64, 1000)
	copy(before, list)

	byteList := UnsafeCastFloatsToBytes(list)
	res := UnsafeCastBytesToFloats(byteList)

	assert.Equal(before, res)

}
