package main

// Command example: gofast -p / -t file -n "filenameOrDirname" -m like

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/0x0FACED/gofast/internal/gofast"
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
	flag.StringVar(&method, "m", "exact", "Search method, e.g., 'exact' or 'pattern'")
	flag.IntVar(&workers, "workers", 10, "Number of goroutines for parallel search (default is 10)")

	log.SetOutput(os.Stdout)
	log.SetFlags(0)
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

	err := gofast.Start(startPath, fileType, name, method, workers)
	if err != nil {
		log.Fatalln(err)
	}

}
