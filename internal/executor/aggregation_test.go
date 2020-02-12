package executor

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"
)

type Register struct {
	registry []int
	n        int
	p        int
	counter  int
}

// n number of items
// p number of groups
func (r *Register) Init(p int, n int) {
	r.n = n
	r.p = p
	r.registry = make([]int, 0)
	for i := 0; i < p; i++ {
		r.registry = append(r.registry, 0)
	}
}

func (r *Register) Update() []int {
	str := big.NewInt(int64(r.counter)).Text(r.n)

	l := []byte(str)
	// fmt.Printf("l %v len %v \n", l, len(l))
	// Poner al reves
	for i := 0; i < len(l); i++ {
		c := r.p - 1 - i
		// fmt.Printf("c %v , i %v \n", c, i)
		x, _ := strconv.Atoi(string(l[len(l)-1-i]))
		r.registry[c] = x

	}
	r.counter++
	fmt.Printf("reg %v \n", r.registry)
	return r.registry
}

func TestHAgg(t *testing.T) {
	//list := setUpTop()
	//a := assert.New(t)
	c := []rune{'a', 'e', 'o', 'r', 'c', 'r', 'd'}
	r := Register{}
	r.Init(6, 7)

	for i := 0; i < 100; i++ {
		r := r.Update()
		for _, it := range r {
			fmt.Printf("%v", string(c[it]))
		}
		fmt.Print(" \n")
	}

}
