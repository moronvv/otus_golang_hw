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

	tasksCnt := 0
	tasksMtx := sync.Mutex{}
	allTasksDone := make(chan struct{})

	errorsExceeded := make(chan struct{})
	errMtx := sync.Mutex{}
	errCnt := 0

	quit := make(chan struct{})

	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case <-quit:
					return
				case task := <-taskCh:
					if err := task(); err != nil {
						func() {
							errMtx.Lock()
							defer errMtx.Unlock()

							select {
							case <-quit:
								return
							default:
								errCnt++
								if errCnt >= m {
									close(quit)
									close(errorsExceeded)
								}
							}
						}()
					}

					func() {
						tasksMtx.Lock()
						defer tasksMtx.Unlock()

						select {
						case <-quit:
							return
						default:
							tasksCnt++
							if tasksCnt >= len(tasks) {
								close(quit)
								close(allTasksDone)
							}
						}
					}()
				}
			}
		}()
	}

	for _, task := range tasks {
		taskCh <- task
	}

	<-quit
	select {
	case <-allTasksDone:
		fmt.Println("success")
		return nil
	case <-errorsExceeded:
		fmt.Println("fatal")
		return ErrErrorsLimitExceeded
	}
}
