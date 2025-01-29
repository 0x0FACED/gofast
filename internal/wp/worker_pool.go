package wp

import (
	"errors"
	"path/filepath"
	"sync"
)

type WorkerPool struct {
	numWorkers   int
	taskCh       <-chan string
	resultCh     chan<- string
	errorCh      chan<- error
	done         chan struct{}
	searchMethod string
	searchName   string
}

func New(numWorkers int, taskCh <-chan string, resultCh chan<- string, errorCh chan<- error, done chan struct{}, method, name string) *WorkerPool {
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
	defer wg.Done()
	var workerWG sync.WaitGroup

	for i := 0; i < wp.numWorkers; i++ {
		workerWG.Add(1)
		go wp.worker(&workerWG)
	}

	workerWG.Wait()
	close(wp.resultCh)
	close(wp.errorCh)
	close(wp.done)
}

func (wp *WorkerPool) worker(wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range wp.taskCh {
		baseName := filepath.Base(task)

		switch wp.searchMethod {
		case "exact":
			if baseName == wp.searchName {
				wp.resultCh <- task
			}
		case "pattern":
			match, err := filepath.Match(wp.searchName, baseName)
			if err != nil {
				wp.errorCh <- err
			} else if match {
				wp.resultCh <- task
			}
		default:
			wp.errorCh <- errors.New("unknown search method: " + wp.searchMethod)
		}
	}
}
