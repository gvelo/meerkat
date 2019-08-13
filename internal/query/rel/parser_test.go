package rel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessQuery(t *testing.T) {

	a := assert.New(t)

	str := "sourcetype=Index campo1=100 and ( campo2=we or campo3=s123 )"
	is := ProcessQuery(str)

	a.Equal("Index", is.IndexScan.indexName)
	a.NotNil(is.IndexScan.filter)

	f := is.IndexScan.filter
	a.NotNil(f.R)
	a.NotNil(f.Op)
	a.NotNil(f.L)

}
