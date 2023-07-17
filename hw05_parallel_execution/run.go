package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskCh := make(chan Task, len(tasks))
	running := sync.WaitGroup{}

	tasksCnt := 0
	errCnt := 0

	taskDone := make(chan struct{})
	taskErr := make(chan struct{})

	quit := make(chan struct{})

	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case <-quit:
					return
				case task, ok := <-taskCh:
					if !ok {
						return
					}

					running.Add(1)

					if err := task(); err != nil {
						taskErr <- struct{}{}
					}

					taskDone <- struct{}{}

					running.Done()
				}
			}
		}()
	}

	for _, task := range tasks {
		taskCh <- task
	}
	close(taskCh)

	for {
		select {
		case <-taskDone:
			tasksCnt++
			if tasksCnt >= len(tasks) {
				fmt.Println("success")
				close(quit)
				return nil
			}
		case <-taskErr:
			errCnt++
			if errCnt >= m {
				fmt.Println("error")
				close(quit)
				// running.Wait()
				return ErrErrorsLimitExceeded
			}
		}
	}
}
