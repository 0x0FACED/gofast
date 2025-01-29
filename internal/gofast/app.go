package gofast

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/0x0FACED/gofast/internal/logs"
	"github.com/0x0FACED/gofast/internal/wp"
)

const (
	DIR  = "dir"
	FILE = "file"
)

type gofast struct {
	wp       *wp.WorkerPool
	taskCh   chan string
	resultCh chan string
	errorCh  chan error
	done     chan struct{}
	start    string
	fileType string
	logger   *logs.Logger
}

func newApp(wp *wp.WorkerPool, t chan string, r chan string, e chan error, d chan struct{}) *gofast {
	l, err := logs.New()
	if err != nil {
		panic(err)
	}
	return &gofast{
		wp:       wp,
		taskCh:   t,
		resultCh: r,
		errorCh:  e,
		done:     d,
		logger:   l,
	}
}

func Start(start, fileType, name, method string, workers int) error {
	taskCh := make(chan string, 500)
	resCh := make(chan string, 500)
	errCh := make(chan error, 50)
	done := make(chan struct{})

	wp := wp.New(workers, taskCh, resCh, errCh, done, method, name)
	app := newApp(wp, taskCh, resCh, errCh, done)
	app.start = start
	app.fileType = fileType

	return app.run()
}

func (a *gofast) run() error {
	var wgApp sync.WaitGroup
	var wgWP sync.WaitGroup

	wgWP.Add(1)
	go a.wp.Start(&wgWP)

	wgApp.Add(1)
	go func() {
		defer wgApp.Done()
		a.processDir(a.start, &wgApp)
	}()

	go a.listen()

	wgApp.Wait()
	close(a.taskCh)

	wgWP.Wait() // Ждём завершения WorkerPool перед завершением программы

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
			a.logger.Info(res)
		case err, ok := <-a.errorCh:
			if !ok {
				a.errorCh = nil
				continue
			}
			a.logger.Error(err)
		case <-a.done:
			return
		}

		if a.resultCh == nil && a.errorCh == nil {
			return
		}
	}
}

func (a *gofast) processDir(path string, wg *sync.WaitGroup) {
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
			}
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				a.processDir(p, wg)
			}(fullPath)
		} else if a.fileType == FILE {
			a.taskCh <- fullPath
		}
	}
}
