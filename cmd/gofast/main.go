package main

// Command example: gofast -p / -t file -n "filenameOrDirname" -m "exact"

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/0x0FACED/gofast/internal/gofast"
	"github.com/0x0FACED/gofast/internal/logs"
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

	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

func main() {
	if os.Args[1] == "logs" {
		f, err := os.Open(logs.LOG_FILENAME)
		if err != nil {
			fmt.Println(err)
			return
		}
		finfo, _ := os.Stat(logs.LOG_FILENAME)
		content := make([]byte, finfo.Size())
		_, err = f.Read(content)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(content))
		return
	}
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

	workers = runtime.NumCPU()

	fmt.Println("Max workers available:", workers)

	err := gofast.Start(startPath, fileType, name, method, workers)
	if err != nil {
		log.Fatalln(err)
	}
}
