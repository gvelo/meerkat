package physical

// type ScanOp struct {
// 	segment storage.Segment
// 	filter  logical.Node
// 	columns []string
// 	bufLen  int
// 	execCtx exec.ExecutionContext
// 	colChs  map[string]chan vector.Vector
// }
//
// func (s *ScanOp) Init() {
// }
//
// func (s *ScanOp) Close() {
// 	// do nothing, segment will be released by the executor
// 	// at the end of the query execution.
// }
//
// func (s *ScanOp) Next() Batch {
//
// 	batch := NewBatch()
//
// 	for name, col := range s.colChs {
//
// 	}
//
// }
