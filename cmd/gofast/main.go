package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	startPath string
	fileType  string
	name      string
	method    string
	workers   int
)

func init() {
	flag.StringVar(&startPath, "p", "/", "Starting point for the search (required parameter)")
	flag.StringVar(&fileType, "t", "file", "Type of search: 'file' or 'dir' (required parameter)")
	flag.StringVar(&name, "n", "", "Name of the file or directory in quotes (required parameter)")
	flag.StringVar(&method, "m", "like", "Search method, e.g., 'like' or 'pattern'")
	flag.IntVar(&workers, "workers", 10, "Number of goroutines for parallel search (default is 10)")
}

func main() {
	flag.Parse()

	if name == "" {
		fmt.Println("Error: parameter -n is required.")
		flag.Usage()
		os.Exit(1)
	}

	if fileType != "file" && fileType != "dir" {
		fmt.Println("Error: parameter -t must be 'file' or 'dir'")
		os.Exit(1)
	}

	results := make(chan string, 200)
	errors := make(chan error, 10)

	var errorList []error
	var errorMu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		startSearch(startPath, fileType, name, method, workers, results, errors)
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	fmt.Println("Search results:")
	go func() {
		for err := range errors {
			errorMu.Lock()
			errorList = append(errorList, err)
			errorMu.Unlock()
		}
	}()

	for {
		select {
		case result, ok := <-results:
			if ok {
				fmt.Println(result)
			} else {
				results = nil
			}
		case <-done:
			if results == nil {
				if len(errorList) > 0 {
					fmt.Println("\nSearch errors:")
					//for _, err := range errorList {
					//	fmt.Fprintf(os.Stderr, "- %v\n", err)
					//}
				}
				return
			}
		}
	}
}

func startSearch(startPath, fileType, name, method string, workers int, results chan<- string, errors chan<- error) {
	tasks := make(chan string)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(tasks, fileType, name, method, results, errors)
		}()
	}

	go func() {
		defer close(tasks)
		// ignore error for now
		_ = filepath.Walk(startPath, func(path string, info os.FileInfo, err error) error {
			tasks <- path
			return nil
		})
	}()

	wg.Wait()
	close(results)
	close(errors)
}

func worker(tasks <-chan string, fileType, name, method string, results chan<- string, errors chan<- error) {
	for path := range tasks {
		info, err := os.Stat(path)
		if err != nil {
			errors <- fmt.Errorf("ошибка чтения информации о %s: %w", path, err)
			continue
		}

		if (fileType == "file" && !info.IsDir()) || (fileType == "dir" && info.IsDir()) {
			switch method {
			case "like":
				if strings.Contains(info.Name(), name) {
					results <- path
				}
			case "pattern":
				matched, err := filepath.Match(name, info.Name())
				if err != nil {
					errors <- fmt.Errorf("ошибка проверки шаблона '%s': %w", name, err)
				} else if matched {
					results <- path
				}
			default:
				errors <- fmt.Errorf("неподдерживаемый метод поиска: %s", method)
			}
		}
	}
}
