package inmem

import (
	"eventdb/segment"
	"fmt"
	"math/rand"
	"testing"
)

func makeTestArray(size int) []segment.Event {
	a := make([]segment.Event, size)

	for i := 0; i < size; i++ {
		a[i] = segment.Event{"key": i & 0xeeeeee, "order": i}
	}

	return a
}

func IsSorted(a []segment.Event, lessThan LessThan) bool {
	len := len(a)

	if len < 2 {
		return true
	}

	prev := a[0]
	for i := 1; i < len; i++ {
		if lessThan(a[i], prev) {
			return false
		}
		prev = a[i]
	}

	return true
}

// use this comparator for sorting
func KeyLessThan(a, b segment.Event) bool {
	return a["key"].(int) < b["key"].(int)
}

type KeyLessThanSlice []segment.Event

func (s KeyLessThanSlice) Len() int {
	return len(s)
}

func (s KeyLessThanSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s KeyLessThanSlice) Less(i, j int) bool {
	return s[i]["key"].(int) < s[j]["key"].(int)
}

// use this comparator to validate sorted data (and prove its stable)
func KeyOrderLessThan(a, b segment.Event) bool {
	if a["key"].(int) < b["key"].(int) {
		return true
	} else if a["key"].(int) == b["key"].(int) {
		return a["order"].(int) < b["order"].(int)
	}

	return false
}

// use this comparator to restore the original order of elements (by sorting on order field)
func OrderLessThan(a, b segment.Event) bool {
	return a["order"].(int) < b["order"].(int)
}

type OrderLessThanSlice []segment.Event

func (s OrderLessThanSlice) Len() int {
	return len(s)
}

func (s OrderLessThanSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s OrderLessThanSlice) Less(i, j int) bool {
	return s[i]["order"].(int) < s[j]["order"].(int)
}

func TestOrderedSameKey(t *testing.T) {
	a := make([]segment.Event, 7)
	a[0] = segment.Event{"key": 3, "order": 4}
	a[1] = segment.Event{"key": 1, "order": 0}
	a[2] = segment.Event{"key": 2, "order": 1}
	a[3] = segment.Event{"key": 2, "order": 2}
	a[4] = segment.Event{"key": 2, "order": 3}
	a[5] = segment.Event{"key": 5, "order": 5}
	a[6] = segment.Event{"key": 6, "order": 6}

	err := Sort(a, KeyLessThan)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}
	t.Logf("key sort: %v", a)

	if !IsSorted(a, OrderLessThan) {
		t.Error("not sorted by same key")
	}

}

func TestNotOrderedSameKey(t *testing.T) {
	a := make([]segment.Event, 7)
	a[0] = segment.Event{"key": 3, "order": 4}
	a[1] = segment.Event{"key": 1, "order": 0}
	a[2] = segment.Event{"key": 2, "order": 1}
	a[3] = segment.Event{"key": 2, "order": 2}
	a[4] = segment.Event{"key": 2, "order": 3}
	a[5] = segment.Event{"key": 5, "order": 5}
	a[6] = segment.Event{"key": 6, "order": 6}
	err := Sort(a, KeyLessThan)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}

	if !IsSorted(a, OrderLessThan) {
		t.Error("not sorted by same key")
	}

}

func TestOrderedSameKey2(t *testing.T) {
	a := make([]segment.Event, 10)
	a[0] = segment.Event{"key": 3, "order": 5}
	a[1] = segment.Event{"key": 1, "order": 0}
	a[2] = segment.Event{"key": 2, "order": 1}
	a[3] = segment.Event{"key": 2, "order": 2}
	a[4] = segment.Event{"key": 2, "order": 3}
	a[5] = segment.Event{"key": 5, "order": 6}
	a[6] = segment.Event{"key": 6, "order": 8}
	a[7] = segment.Event{"key": 2, "order": 4}
	a[8] = segment.Event{"key": 5, "order": 7}
	a[9] = segment.Event{"key": 6, "order": 9}

	err := Sort(a, KeyLessThan)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}

	if !IsSorted(a, OrderLessThan) {
		t.Error("not sorted by same key")
	}

}

func TestSmoke(t *testing.T) {
	a := make([]segment.Event, 3)
	a[0] = segment.Event{"key": 3, "order": 0}
	a[1] = segment.Event{"key": 1, "order": 1}
	a[2] = segment.Event{"key": 2, "order": 2}

	err := Sort(a, KeyLessThan)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}

	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func TestSmokeTS(t *testing.T) {
	a := make([]segment.Event, 3)
	a[0] = segment.Event{"key": 3, "order": 0}
	a[1] = segment.Event{"key": 1, "order": 1}
	a[2] = segment.Event{"key": 2, "order": 2}

	err := TimSort(KeyLessThanSlice(a))
	if err != nil {
		t.Fatalf("sort: %v", err)
	}

	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func TestSmokeStability(t *testing.T) {
	a := make([]segment.Event, 3)
	a[0] = segment.Event{"key": 3, "order": 0}
	a[1] = segment.Event{"key": 2, "order": 1}
	a[2] = segment.Event{"key": 2, "order": 2}

	err := Sort(a, KeyLessThan)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}

	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func TestSmokeStabilityTS(t *testing.T) {
	a := make([]segment.Event, 3)
	a[0] = segment.Event{"key": 3, "order": 0}
	a[1] = segment.Event{"key": 2, "order": 1}
	a[2] = segment.Event{"key": 2, "order": 2}

	err := TimSort(KeyLessThanSlice(a))
	if err != nil {
		t.Fatalf("sort: %v", err)
	}

	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func Test1K(t *testing.T) {
	a := makeTestArray(1024)

	err := Sort(a, KeyLessThan)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func Test1KTS(t *testing.T) {
	a := makeTestArray(1024)

	err := TimSort(KeyLessThanSlice(a))
	if err != nil {
		t.Fatalf("sort: %v", err)
	}
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func Test1M(t *testing.T) {
	a := makeTestArray(1024 * 1024)

	err := Sort(a, KeyLessThan)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func Test1MTS(t *testing.T) {
	a := makeTestArray(1024 * 1024)

	err := TimSort(KeyLessThanSlice(a))
	if err != nil {
		t.Fatalf("sort: %v", err)
	}
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func makeRandomArray(size int) []segment.Event {
	a := make([]segment.Event, size)

	for i := 0; i < size; i++ {
		a[i] = segment.Event{"key": rand.Intn(100), "order": i}
	}

	return a
}

func Equals(a, b segment.Event) bool {
	return a["key"].(int) == b["key"].(int) && a["order"].(int) == b["order"].(int)
}

func TestRandom1M(t *testing.T) {
	size := 1024 * 1024

	a := makeRandomArray(size)
	b := make([]segment.Event, size)
	copy(b, a)

	err := Sort(a, KeyLessThan)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}

	// sort by order
	err = Sort(a, OrderLessThan)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}
	for i := 0; i < len(b); i++ {
		if !Equals(b[i], a[i]) {
			t.Error("oops")
		}
	}
}

func TestRandom1MTS(t *testing.T) {
	size := 1024 * 1024

	a := makeRandomArray(size)
	b := make([]segment.Event, size)
	copy(b, a)

	err := TimSort(KeyLessThanSlice(a))
	if err != nil {
		t.Fatalf("sort: %v", err)
	}
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}

	// sort by order
	err = TimSort(OrderLessThanSlice(a))
	if err != nil {
		t.Fatalf("sort: %v", err)
	}
	for i := 0; i < len(b); i++ {
		if !Equals(b[i], a[i]) {
			t.Error("oops")
		}
	}
}

type person struct {
	ssn  int
	name string
}

func (p *person) Less(o interface{}) bool {
	return p.ssn < o.(*person).ssn
}

const (
	_Sawtooth = iota
	_Rand
	_Stagger
	_Plateau
	_Shuffle
	_NDist
)

const (
	_Copy = iota
	_Reverse
	_ReverseFirstHalf
	_ReverseSecondHalf
	_Sorted
	_Dither
	_NMode
)

func TestBentleyMcIlroy(t *testing.T) {
	//	sizes := []int{100, 1023, 1024, 1025, 1023 * 1023, 1024 * 1024, 1025 * 1025}
	sizes := []int{100, 1023, 1024, 1025}
	dists := []string{"sawtooth", "rand", "stagger", "plateau", "shuffle"}
	modes := []string{"copy", "reverse", "reverse1", "reverse2", "sort", "dither"}
	tmp1 := make([]segment.Event, 1025*1025)
	tmp2 := make([]segment.Event, 1025*1025)
	for ni := 0; ni < len(sizes); ni++ {
		n := sizes[ni]
		for m := 1; m < 2*n; m *= 2 {
			for dist := 0; dist < _NDist; dist++ {
				j := 0
				k := 1
				data := tmp1[0:n]
				for i := 0; i < n; i++ {
					switch dist {
					case _Sawtooth:
						data[i] = segment.Event{"key": i % m, "order": i}
					case _Rand:
						data[i] = segment.Event{"key": rand.Intn(m), "order": i}
					case _Stagger:
						data[i] = segment.Event{"key": (i*m + i) % n, "order": i}
					case _Plateau:
						if i < m {
							data[i] = segment.Event{"key": i, "order": i}
						} else {
							data[i] = segment.Event{"key": m, "order": i}
						}
					case _Shuffle:
						if rand.Intn(m) != 0 {
							j += 2
							data[i] = segment.Event{"key": j, "order": i}
						} else {
							k += 2
							data[i] = segment.Event{"key": k, "order": i}
						}
					}
				}

				mdata := tmp2[0:n]
				for mode := 0; mode < _NMode; mode++ {
					switch mode {
					case _Copy:
						for i := 0; i < n; i++ {
							mdata[i] = segment.Event{"key": data[i]["key"].(int), "order": i}
						}
					case _Reverse:
						for i := 0; i < n; i++ {
							mdata[i] = segment.Event{"key": data[n-i-1]["key"].(int), "order": i}
						}
					case _ReverseFirstHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = segment.Event{"key": data[n/2-i-1]["key"].(int), "order": i}
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = segment.Event{"key": data[i]["key"].(int), "order": i}
						}
					case _ReverseSecondHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = segment.Event{"key": data[i]["key"].(int), "order": i}
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = segment.Event{"key": data[n-(i-n/2)-1]["key"].(int), "order": i}
						}
					case _Sorted:
						for i := 0; i < n; i++ {
							mdata[i] = data[i]
						}
						// SortInts is known to be correct
						// because mode Sort runs after mode _Copy.
						Sort(mdata, KeyLessThan)
					case _Dither:
						for i := 0; i < n; i++ {
							mdata[i] = segment.Event{"key": data[i]["key"].(int) + i%5, "order": i}
						}
					}

					desc := fmt.Sprintf("n=%d m=%d dist=%s mode=%s", n, m, dists[dist], modes[mode])

					for i := 0; i < len(mdata); i++ {
						mdata[i] = segment.Event{"key": mdata[i]["key"].(int), "order": i}
					}

					gdata := make([]segment.Event, len(mdata))
					copy(gdata, mdata)

					err := Sort(mdata, KeyLessThan)
					if err != nil {
						t.Fatalf("sort: %v", err)
					}

					// If we were testing C qsort, we'd have to make a copy
					// of the array and sort it ourselves and then compare
					// x against it, to ensure that qsort was only permuting
					// the data, not (for example) overwriting it with zeros.
					//
					// In go, we don't have to be so paranoid: since the only
					// mutating method Sort can call is TestingData.swap,
					// it suffices here just to check that the final array is sorted.
					if !IsSorted(mdata, KeyOrderLessThan) {
						t.Errorf("%s: ints not sorted", desc)
						t.Errorf("\t%v", mdata)
						t.FailNow()
					}

					err = Sort(mdata, OrderLessThan)
					if err != nil {
						t.Fatalf("sort: %v", err)
					}
					for i := 0; i < len(data); i++ {
						if !Equals(gdata[i], mdata[i]) {
							t.Error("restore sort failed")
							t.Errorf("gdata=%v", gdata)
							t.Errorf("mdata=%v", mdata)
							t.Errorf("bad index: %v\n", i)
							t.FailNow()
						}
					}
				}
			}
		}
	}
}

func TestBentleyMcIlroyTS(t *testing.T) {
	//	sizes := []int{100, 1023, 1024, 1025, 1023 * 1023, 1024 * 1024, 1025 * 1025}
	sizes := []int{100, 1023, 1024, 1025}
	dists := []string{"sawtooth", "rand", "stagger", "plateau", "shuffle"}
	modes := []string{"copy", "reverse", "reverse1", "reverse2", "sort", "dither"}
	tmp1 := make([]segment.Event, 1025*1025)
	tmp2 := make([]segment.Event, 1025*1025)
	for ni := 0; ni < len(sizes); ni++ {
		n := sizes[ni]
		for m := 1; m < 2*n; m *= 2 {
			for dist := 0; dist < _NDist; dist++ {
				j := 0
				k := 1
				data := tmp1[0:n]
				for i := 0; i < n; i++ {
					switch dist {
					case _Sawtooth:
						data[i] = segment.Event{"key": i % m, "order": i}
					case _Rand:
						data[i] = segment.Event{"key": rand.Intn(m), "order": i}
					case _Stagger:
						data[i] = segment.Event{"key": (i*m + i) % n, "order": i}
					case _Plateau:
						if i < m {
							data[i] = segment.Event{"key": i, "order": i}
						} else {
							data[i] = segment.Event{"key": m, "order": i}
						}
					case _Shuffle:
						if rand.Intn(m) != 0 {
							j += 2
							data[i] = segment.Event{"key": j, "order": i}
						} else {
							k += 2
							data[i] = segment.Event{"key": k, "order": i}
						}
					}
				}

				mdata := tmp2[0:n]
				for mode := 0; mode < _NMode; mode++ {
					switch mode {
					case _Copy:
						for i := 0; i < n; i++ {
							mdata[i] = segment.Event{"key": data[i]["key"].(int), "order": i}
						}
					case _Reverse:
						for i := 0; i < n; i++ {
							mdata[i] = segment.Event{"key": data[n-i-1]["key"].(int), "order": i}
						}
					case _ReverseFirstHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = segment.Event{"key": data[n/2-i-1]["key"].(int), "order": i}
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = segment.Event{"key": data[i]["key"].(int), "order": i}
						}
					case _ReverseSecondHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = segment.Event{"key": data[i]["key"].(int), "order": i}
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = segment.Event{"key": data[n-(i-n/2)-1]["key"].(int), "order": i}
						}
					case _Sorted:
						for i := 0; i < n; i++ {
							mdata[i] = data[i]
						}
						// SortInts is known to be correct
						// because mode Sort runs after mode _Copy.
						Sort(mdata, KeyLessThan)
					case _Dither:
						for i := 0; i < n; i++ {
							mdata[i] = segment.Event{"key": data[i]["key"].(int) + i%5, "order": i}
						}
					}

					desc := fmt.Sprintf("n=%d m=%d dist=%s mode=%s", n, m, dists[dist], modes[mode])

					for i := 0; i < len(mdata); i++ {
						mdata[i] = segment.Event{"key": mdata[i]["key"].(int), "order": i}
					}

					gdata := make([]segment.Event, len(mdata))
					copy(gdata, mdata)

					err := TimSort(KeyLessThanSlice(mdata))
					if err != nil {
						t.Fatalf("sort: %v", err)
					}

					// If we were testing C qsort, we'd have to make a copy
					// of the array and sort it ourselves and then compare
					// x against it, to ensure that qsort was only permuting
					// the data, not (for example) overwriting it with zeros.
					//
					// In go, we don't have to be so paranoid: since the only
					// mutating method Sort can call is TestingData.swap,
					// it suffices here just to check that the final array is sorted.
					if !IsSorted(mdata, KeyOrderLessThan) {
						t.Errorf("%s: ints not sorted", desc)
						t.Errorf("\t%v", mdata)
						t.FailNow()
					}

					err = TimSort(OrderLessThanSlice(mdata))
					if err != nil {
						t.Fatalf("sort: %v", err)
					}
					for i := 0; i < len(data); i++ {
						if !Equals(gdata[i], mdata[i]) {
							t.Error("restore sort failed")
							t.Errorf("gdata=%v", gdata)
							t.Errorf("mdata=%v", mdata)
							t.Errorf("bad index: %v\n", i)
							t.FailNow()
						}
					}
				}
			}
		}
	}
}
