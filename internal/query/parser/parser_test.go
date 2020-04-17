package parser

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {

	s := ` StormEvents
          | limit 10
          | count`

	p := NewParser()
	stmt, err := p.Parse(s)

	if err != nil {
		panic(err)
	}

	fmt.Println(stmt)

}
