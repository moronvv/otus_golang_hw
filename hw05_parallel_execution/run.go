package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task)

	errCnt := 0
	var errMtx sync.Mutex

	var workersWg sync.WaitGroup
	workersWg.Add(n)

	for i := 0; i < n; i++ {
		// workers
		go func() {
			defer workersWg.Done()

			for task := range tasksCh {
				if err := task(); err != nil {
					errMtx.Lock()
					errCnt++
					errMtx.Unlock()
				}
			}
		}()
	}

	for _, task := range tasks {
		errMtx.Lock()
		errorsExceeded := errCnt >= m
		errMtx.Unlock()
		if errorsExceeded {
			break
		}

		tasksCh <- task
	}
	close(tasksCh)

	workersWg.Wait()

	if errCnt >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
