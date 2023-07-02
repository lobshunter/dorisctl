package step

import (
	"context"
	"runtime"

	"golang.org/x/sync/errgroup"

	"github.com/lobshunter/dorisctl/pkg/cluster/task"
)

var (
	_ Step = &Serial{}
	_ Step = &Parallel{}
)

type Step interface {
	task.Task
	StepName() string // StepName is used to seprate Task and Step type
}

type Serial struct {
	inner []task.Task
}

type Parallel struct {
	nprocs int // TODO: configurable limit number of goroutines
	inner  []task.Task
}

func NewSerial(inner ...task.Task) *Serial {
	return &Serial{
		inner: inner,
	}
}

func (s *Serial) Name() string {
	return "Serial"
}

func (s *Serial) StepName() string {
	return s.Name()
}

func (s *Serial) Execute(ctx context.Context) error {
	var err error
	for _, t := range s.inner {
		if err = t.Execute(ctx); err != nil {
			break
		}
	}
	return err
}

func NewParallel(inner ...task.Task) *Parallel {
	return &Parallel{
		nprocs: runtime.NumCPU(),
		inner:  inner,
	}
}

func (p *Parallel) Name() string {
	return "Parallel"
}

func (p *Parallel) StepName() string {
	return p.Name()
}

func (p *Parallel) Execute(ctx context.Context) error {
	errg := errgroup.Group{}
	errg.SetLimit(p.nprocs)

	for _, t := range p.inner {
		tsk := t
		errg.Go(func() error {
			return tsk.Execute(ctx)
		})
	}

	return errg.Wait()
}
