package internal

import (
	"sync"
	"sync/atomic"

	"github.com/dndev-xx/go-os-std/hw05_parallel_execution/internal/view" //nolint:depguard
)

type Task func() error

type Runner struct {
	IRunner
	wg            sync.WaitGroup
	taskChan      chan Task
	executionChan chan error
	signal        chan struct{}
	errorCnt      int32
}

func NewRunner() *Runner {
	return &Runner{
		taskChan:      make(chan Task),
		executionChan: make(chan error),
		signal:        make(chan struct{}),
		errorCnt:      view.ZERO,
	}
}

func (r *Runner) Run(tasks []Task, n, m int) error {
	should, err := r.ShouldExecutionTasks(n, m)
	if err != nil {
		return err
	}
	r.wg.Add(view.UNIT)
	go func() {
		defer r.wg.Done()
		defer close(r.executionChan)
		wg := sync.WaitGroup{}
		wg.Add(n)
		for i := view.ZERO; i < n; i++ {
			go func() {
				defer wg.Done()
				r.StartWorkers()
			}()
		}
		wg.Wait()
	}()
	r.wg.Add(view.UNIT)
	go func() {
		defer r.wg.Done()
		defer close(r.taskChan)
		r.SendTask(tasks)
	}()
	for {
		rsl, ok := <-r.executionChan
		if !ok {
			break
		}
		if *should {
			continue
		}
		if rsl != nil {
			atomic.AddInt32(&r.errorCnt, view.UNIT)
		}
		if int(r.errorCnt) >= m {
			err = view.ErrErrorsLimitExceeded
			close(r.signal)
			break
		}
	}
	r.wg.Wait()
	return err
}

func (r *Runner) StartWorkers() {
	for {
		task, ok := <-r.taskChan
		if !ok {
			return
		}
		select {
		case r.executionChan <- task():
		case <-r.signal:
			return
		}
	}
}

func (r *Runner) SendTask(tasks []Task) {
	for _, task := range tasks {
		select {
		case r.taskChan <- task:
		case <-r.signal:
			return
		}
	}
}

func (r *Runner) ShouldExecutionTasks(n, m int) (*bool, error) {
	if n < view.UNIT {
		return nil, view.ErrErrorsLessUnitWorker
	}
	resolution := m < view.UNIT
	return &resolution, nil
}
