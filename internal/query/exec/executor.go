package exec

import (
	"github.com/google/uuid"
	"sync"
)

//go:generate protoc -I . -I ../../../build/proto/ -I ../../../internal/storage/ --plugin ../../../build/protoc-gen-gogofaster --gogofaster_out=plugins=grpc,paths=source_relative:.  ./exec.proto

type ExecutionContext interface {
	Cancel(err error)
	Done() <-chan struct{}
	Error() error
}

func NewExecutionContext() ExecutionContext {
	return &executionContext{
		done: make(<-chan struct{}),
	}
}

type executionContext struct {
	mu   sync.Mutex
	err  error
	done chan struct{}
}

func (c *executionContext) Cancel(err error) {

	defer c.mu.Unlock()

	c.mu.Lock()

	if c.err != nil {
		return
	}

	c.err = err

	close(c.done)

}

func (c *executionContext) Done() <-chan struct{} {
	return c.done
}

func (c *executionContext) Error() error {
	defer c.mu.Unlock()
	c.mu.Lock()
	return c.err
}

type QueryExecution interface {
	Id() uuid.UUID
	Out() <-chan interface{}
	Cancel(err error)
	Wait()
}

type Executor interface {
	ExecuteQuery(query string) QueryExecution
	CancelQuery(queryId uuid.UUID)
	Stop()
}

func NewExecutor() Executor {
	return &executor{}
}

type executor struct {
}

func (e executor) ExecuteQuery(query string) (QueryExecution, error) {

	coordinator := NewCoordinatorExecutor()

	err := coordinator.exec(query)

	if err != nil {
		return nil, err
	}

	// TODO: keep track of the query execution so we can cancel or
	// wait on it

	return coordinator

}

func (e executor) CancelQuery(queryId uuid.UUID) {
	panic("implement me")
}

func (e executor) Stop() {
	panic("implement me")
}
