package internal

import (
	"sync"
	"sync/atomic"

	"github.com/dndev-xx/go-os-std/hw05_parallel_execution/internal/view" //nolint:depguard
)

type Task func() error

func Run(tasks []Task, n, m int) error { //nolint:gocognit
	var errCnt int32
	var flag atomic.Bool
	executionChan := make(chan Task, len(tasks))
	var wg sync.WaitGroup
	for i := view.ZERO; i < n; i++ {
		wg.Add(view.UNIT)
		go func() {
			defer wg.Done()
			for task := range executionChan {
				if flag.Load() {
					return
				}
				if err := task(); m > view.ZERO {
					if err != nil {
						atomic.AddInt32(&errCnt, view.UNIT)
						if int(errCnt) >= m {
							flag.Store(true)
							return
						}
					}
				}
			}
		}()
	}
	go func() {
		defer close(executionChan)
		for _, task := range tasks {
			if m > 0 && int(errCnt) >= m {
				return
			}
			executionChan <- task
		}
	}()
	wg.Wait()
	if m > view.ZERO && int(errCnt) >= m {
		return view.ErrErrorsLimitExceeded
	}
	return nil
}
