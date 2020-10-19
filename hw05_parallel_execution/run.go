package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
// M <= 0 - ignore errors.
func Run(tasks []Task, n int, m int) error {
	wg := &sync.WaitGroup{}
	taskChan := make(chan Task)
	resultChan := make(chan error, n)

	executeTasks(wg, resultChan, taskChan, n)
	go waitExecuteTasksDone(wg, resultChan)
	isErrorsLimitExceeded := runTasks(tasks, resultChan, taskChan, m)

	close(taskChan)
	for range resultChan {
	}

	if isErrorsLimitExceeded {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func executeTasks(wg *sync.WaitGroup, resultChan chan<- error, taskChan <-chan Task, n int) {
	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for task := range taskChan {
				resultChan <- task()
			}
		}()
	}
}

func waitExecuteTasksDone(wg *sync.WaitGroup, resultChan chan error) {
	wg.Wait()
	close(resultChan)
}

func runTasks(tasks []Task, resultChan <-chan error, taskChan chan<- Task, m int) bool {
	errorsCount := 0

	for i := 0; i < len(tasks); {
		task := tasks[i]

		select {
		case err := <-resultChan:
			if m > 0 && err != nil {
				errorsCount++
				if errorsCount == m {
					return true
				}
			}
		case taskChan <- task:
			i++
		}
	}

	return false
}
