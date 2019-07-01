package utils

type Slicer interface {
	Add(interface{})
	Get() interface{}
	Reset()
	Total() int
}

type SlicerInt struct {
	data []int
}

func (s *SlicerInt) Add(a interface{}) {
	s.data = append(s.data, a.(int))
}

func (s *SlicerInt) Get() interface{} {
	return s.data
}

func (s *SlicerInt) Total() int {
	return len(s.data)
}

func (s *SlicerInt) Reset() {
	s.data = make([]int, 0)
}

type SlicerFloat struct {
	data []float64
}

func (s *SlicerFloat) Add(a interface{}) {
	s.data = append(s.data, a.(float64))
}

func (s *SlicerFloat) Get() interface{} {
	return s.data
}

func (s *SlicerFloat) Reset() {
	s.data = make([]float64, 0)
}

func (s *SlicerFloat) Total() int {
	return len(s.data)
}

type SlicerString struct {
	data []string
}

func (s *SlicerString) Add(a interface{}) {
	s.data = append(s.data, a.(string))
}

func (s *SlicerString) Get() interface{} {
	return s.data
}

func (s *SlicerString) Reset() {
	s.data = make([]string, 0)
}

func (s *SlicerString) Total() int {
	return len(s.data)
}
