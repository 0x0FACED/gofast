package gofast

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/0x0FACED/gofast/configs"
	"github.com/0x0FACED/gofast/internal/wp"
)

const (
	DIR  = "dir"
	FILE = "file"
)

type gofast struct {
	wp *wp.WorkerPool

	taskCh   chan<- string // write to
	resultCh chan string   // read from
	errorCh  chan error    // read from
	done     chan struct{} // read from

	start    string
	fileType string

	cfg configs.AppConfig
}

func new(wp *wp.WorkerPool, t chan string, r chan string, e chan error, d chan struct{}) *gofast {
	return &gofast{
		wp:       wp,
		resultCh: r,
		taskCh:   t,
		errorCh:  e,
		done:     d,
	}
}

func Start(start, fileType, name, method string, workers int) error {
	cfg, err := configs.Load()
	if err != nil {
		return err
	}

	taskCh := make(chan string, cfg.ResultChanSize)
	resCh := make(chan string, cfg.ResultChanSize)
	errCh := make(chan error, cfg.ErrorsChanSize)
	done := make(chan struct{})

	wp := wp.New(workers, taskCh, resCh, errCh, done, method, name)

	app := new(wp, taskCh, resCh, errCh, done)
	app.start = start
	app.fileType = fileType

	return app.run()
}

func (a *gofast) run() error {
	var wgApp sync.WaitGroup
	var wgWP sync.WaitGroup

	go a.wp.Start(&wgWP)

	wgApp.Add(1)
	go func() {
		defer wgApp.Done()
		a.processDir(a.start, &wgApp)
	}()

	go a.listen()

	wgApp.Wait()

	return nil
}

func (a *gofast) listen() {
	for {
		select {
		case res, ok := <-a.resultCh:
			if !ok {
				a.resultCh = nil
				continue
			}
			log.Println(res)
		case _, ok := <-a.errorCh: // ignore
			if !ok {
				a.errorCh = nil
				continue
			}
		case <-a.done:
			log.Println("Worker finished")
		}

		if a.resultCh == nil && a.errorCh == nil {
			return
		}
	}
}

func (a *gofast) processDir(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	entries, err := os.ReadDir(path)
	if err != nil {
		a.errorCh <- err
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			if a.fileType == DIR {
				a.taskCh <- fullPath
				continue
			}
			wg.Add(1)
			go a.processDir(fullPath, wg)
		} else if a.fileType == FILE {
			a.taskCh <- fullPath
		}
	}
}
