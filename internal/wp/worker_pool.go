package wp

import (
	"errors"
	"path/filepath"
	"sync"
)

type WorkerPool struct {
	numWorkers int

	taskCh   <-chan string   // Tasks will be here
	resultCh chan<- string   // after complete task send res to this chan
	errorCh  chan<- error    // if err -> send to errorCh
	done     chan<- struct{} // if there are no tasks -> send signal

	searchMethod string
	searchName   string
}

func New(numWorkers int, taskCh <-chan string, resultCh chan<- string, errorCh chan<- error, done chan<- struct{}, method, name string) *WorkerPool {
	return &WorkerPool{
		numWorkers:   numWorkers,
		taskCh:       taskCh,
		resultCh:     resultCh,
		errorCh:      errorCh,
		done:         done,
		searchMethod: method,
		searchName:   name,
	}
}

func (wp *WorkerPool) Start(wg *sync.WaitGroup) {
	for i := 0; i < wp.numWorkers; i++ {
		wg.Add(1)
		go wp.worker(wg)
	}

	go func() {
		wg.Wait()
		wp.done <- struct{}{}
		close(wp.resultCh)
		close(wp.errorCh)
		close(wp.done)
	}()
}

func (wp *WorkerPool) worker(wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range wp.taskCh {
		// Extract the file/directory name from the path
		baseName := filepath.Base(task)

		switch wp.searchMethod {
		case "exact":
			if baseName == wp.searchName {
				wp.resultCh <- task // Add to results if there's an exact match
			}
		case "pattern":
			match, err := filepath.Match(wp.searchName, baseName)
			if err != nil {
				wp.errorCh <- err // Handle invalid patterns
			} else if match {
				wp.resultCh <- task // Add to results if the pattern matches
			}
		default:
			wp.errorCh <- errors.New("unknown search method: " + wp.searchMethod)
			panic(-1)
		}
	}
}
