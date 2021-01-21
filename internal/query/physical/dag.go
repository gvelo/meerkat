package physical

import (
	"fmt"
	"github.com/google/uuid"
	"meerkat/internal/query/execbase"
	"meerkat/internal/query/execpb"
	"meerkat/internal/storage"
	"runtime/debug"
	"strings"
	"sync"
)

type DAG interface {
	Run()
}

var _ DAG = &executableDAG{}

// the root of the DAG. Handle operator lifecycle and acquired resources.
type executableDAG struct {
	execCtx       execbase.ExecutionContext
	runnables     []RunnableOp
	roots         []RunnableOp
	queryId       uuid.UUID
	wg            *sync.WaitGroup
	segments      []storage.Segment
	segReg        storage.SegmentRegistry
	localNodeName string
}

func NewDAG(
	execCtx execbase.ExecutionContext,
	runnables []RunnableOp,
	roots []RunnableOp,
	queryId uuid.UUID,
	segments []storage.Segment,
	segReg storage.SegmentRegistry,
	localNodeName string,
) *executableDAG {

	return &executableDAG{
		execCtx:       execCtx,
		runnables:     runnables,
		roots:         roots,
		queryId:       queryId,
		wg:            &sync.WaitGroup{},
		segments:      segments,
		segReg:        segReg,
		localNodeName: localNodeName,
	}

}

func (ed *executableDAG) Run() {

	// TODO(gvelo) init operators

	ed.Dump()

	for _, runnable := range ed.runnables {
		ed.wg.Add(1)
		fmt.Printf("running op %T \n", runnable)
		go ed.runOp(runnable)
	}

	ed.wg.Wait()

	// TODO(gvelo) close operators

	ed.releaseSegments()

}

func (ed *executableDAG) runOp(op RunnableOp) {

	defer ed.wg.Done()

	defer func() {
		if r := recover(); r != nil {

			debug.PrintStack()

			if execErr, ok := r.(*execpb.ExecError); ok {
				ed.execCtx.CancelWithExecError(execErr)
				return
			}

			execErr := execbase.ExtractExecError(r)

			if execErr != nil {
				ed.execCtx.CancelWithExecError(execErr)
				return
			}

			execErr = execbase.NewExecError(
				fmt.Sprintf("Error executing query: %v", r),
				ed.localNodeName,
			)

			ed.execCtx.CancelWithExecError(execErr)

		}
	}()

	op.Run()
}

func (ed *executableDAG) releaseSegments() {

	for _, segment := range ed.segments {
		id, err := uuid.FromBytes(segment.Info().Id)
		if err != nil {
			panic(err)
		}
		ed.segReg.Release(id)
	}

}

func (ed *executableDAG) Dump() {
	v := &dagDumpVisitor{}
	Walk(ed.roots[0], v)
}

type dagDumpVisitor struct {
	ident int
}

func (d *dagDumpVisitor) VisitPre(n Operator) Operator {
	d.ident += 10
	sp := strings.Repeat(" ", d.ident)
	fmt.Printf("%v %T \n", sp, n)
	return n
}

func (d *dagDumpVisitor) VisitPost(n Operator) Operator {
	d.ident -= 10
	return n
}
