package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var taskIndex int32
	var errorCount int32
	taskCount := int32(len(tasks))
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				atomic.AddInt32(&taskIndex, 1)
				if atomic.LoadInt32(&taskIndex)-1 >= taskCount {
					return
				}

				if atomic.LoadInt32(&errorCount) > int32(m) && m > 0 {
					return
				}

				task := tasks[atomic.LoadInt32(&taskIndex)-1]
				if task == nil {
					continue
				}

				err := task()

				if err != nil && m > 0 {
					atomic.AddInt32(&errorCount, 1)
				}
			}
		}()
	}
	wg.Wait()
	if atomic.LoadInt32(&errorCount) > 0 && m > 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
