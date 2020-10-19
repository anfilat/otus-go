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
	done := make(chan interface{})
	taskChan := make(chan Task)
	resultChan := make(chan error, n)

	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(wg, done, resultChan, taskChan)
	}

	isErrorsLimitExceeded := runTasks(tasks, resultChan, taskChan, m)

	close(done)
	close(taskChan)
	wg.Wait()

	if isErrorsLimitExceeded {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(wg *sync.WaitGroup, done <-chan interface{}, resultChan chan<- error, taskChan <-chan Task) {
	defer wg.Done()

	for task := range taskChan {
		result := task()
		select {
		case <-done:
			return
		case resultChan <- result:
		}
	}
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
