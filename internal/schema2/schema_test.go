package schema2

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func TestDbSerde(t *testing.T) {

	s := `{"name":"testname"}`

	a := &DatabaseDesc{}
//	fmt.Println(a.Name())
	err := jsoniter.Unmarshal([]byte(s), a)
	fmt.Println(err)
	fmt.Println(a.Name())

}
