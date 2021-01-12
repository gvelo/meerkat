package physical

import (
	"fmt"
	"github.com/google/uuid"
	"meerkat/internal/query/exec"
	"meerkat/internal/storage"
	"runtime/debug"
	"sync"
)

type DAG interface {
	Run()
}

var _ DAG = &executableDAG{}

// the root of the DAG. Handle operator lifecicle and acquired resources.
type executableDAG struct {
	execCtx   exec.ExecutionContext
	runnables []RunnableOp
	roots     []RunnableOp
	queryId   uuid.UUID
	wg        *sync.WaitGroup
	segments  []storage.Segment
	segReg    storage.SegmentRegistry
}

func NewDAG(
	execCtx exec.ExecutionContext,
	runnables []RunnableOp,
	roots []RunnableOp,
	queryId uuid.UUID,
	segments []storage.Segment,
	segReg storage.SegmentRegistry,
) *executableDAG {

	return &executableDAG{
		execCtx:   execCtx,
		runnables: runnables,
		roots:     roots,
		queryId:   queryId,
		wg:        &sync.WaitGroup{},
		segments:  segments,
		segReg:    segReg,
	}

}

func (ed *executableDAG) Run() {

	// TODO(gvelo) init operators

	for _, runnable := range ed.runnables {
		ed.wg.Add(1)
		go ed.runOp(runnable)
	}

	wg.Wait()

	// TODO(gvelo) close operators

	ed.releaseSegments()

}

func (ed *executableDAG) runOp(op RunnableOp) {

	defer ed.wg.Done()

	defer func() {
		if r := recover(); r != nil {
			// TODO(gvelo) log propertly
			// handle all errors here .. ( check remote error etc .. )
			// no need to pass a cancel func to all the gorutines
			//
			debug.PrintStack()
			error := fmt.Errorf("error executing query : %v", r)
			ed.execCtx.CancelWithPropagation(error, &exec.ExecError{
				Id:     ed.queryId[:],
				Detail: fmt.Sprintf("error on operator %T", ed),
			})
		}
	}()

	ed.Run()
}

func (ed executableDAG) releaseSegments() {

	for _, segment := range ed.segments {
		id, err := uuid.FromBytes(segment.Info().Id)
		if err != nil {
			panic(err)
		}
		ed.segReg.Release(id)
	}

}
