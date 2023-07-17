package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(tasks chan Task, wg *sync.WaitGroup, allTasksDone chan struct{}, taskErr chan error, quit chan struct{}) {
	defer wg.Done()

	for {
		select {
		case <-quit:
			return
		case task, ok := <-tasks:
			if !ok {
				allTasksDone <- struct{}{}
				return
			}

			if err := task(); err != nil {
				taskErr <- err
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task, len(tasks))
	var result error

	taskErr := make(chan error)
	allTasksDone := make(chan struct{})
	quit := make(chan struct{})
	var workersWg sync.WaitGroup

	errCnt := 0

	for i := 0; i < n; i++ {
		workersWg.Add(1)
		go worker(tasksCh, &workersWg, allTasksDone, taskErr, quit)
	}

	for _, task := range tasks {
		tasksCh <- task
	}

loop:
	for {
		select {
		case <-taskErr:
			errCnt++
			if errCnt >= m {
				close(quit)
				close(tasksCh)
				fmt.Println("error")
				result = ErrErrorsLimitExceeded
				break loop
			}
		case <-allTasksDone:
			close(tasksCh)
			fmt.Println("done")
			break loop
		}
	}

	fmt.Println("wait all workers done")
	workersWg.Wait()

	return result
}
