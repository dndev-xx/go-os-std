package internal

import (
	"sync"

	"github.com/dndev-xx/go-os-std/hw05_parallel_execution/internal/view" //nolint:depguard
)

type Task func() error

func Run(tasks []Task, n, m int) error {
	if n < view.UNIT {
		return view.ErrErrorsLessUnitWorker
	}
	resolution := m < view.UNIT
	taskChan := make(chan Task)
	executionChan := make(chan error)
	signal := make(chan struct{})
	var errorCnt int32 = view.ZERO
	var outerWg sync.WaitGroup
	outerWg.Add(view.UNIT)
	go func() {
		defer outerWg.Done()
		defer close(executionChan)
		wg := sync.WaitGroup{}
		wg.Add(n)
		for i := view.ZERO; i < n; i++ {
			go func() {
				defer wg.Done()
				for task := range taskChan {
					select {
					case executionChan <- task():
					case <-signal:
					}
				}
			}()
		}
		wg.Wait()
	}()
	outerWg.Add(view.UNIT)
	go func() {
		defer outerWg.Done()
		defer close(taskChan)
		for _, task := range tasks {
			select {
			case taskChan <- task:
			case <-signal:
				return
			}
		}
	}()
	var err error
	for {
		rsl, ok := <-executionChan
		if !ok {
			break
		}
		if resolution {
			continue
		}
		if rsl != nil {
			errorCnt++
		}
		if int(errorCnt) >= m {
			err = view.ErrErrorsLimitExceeded
			close(signal)
			break
		}
	}
	outerWg.Wait()
	return err
}
