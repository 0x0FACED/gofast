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

const (
	VERSION = "1.1.2"
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
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "logs":
			if len(os.Args) > 2 && os.Args[2] == "clear" {
				err := os.Truncate(logs.LOG_FILENAME, 0)
				if err != nil {
					fmt.Println("Error clearing logs:", err)
					return
				}
				fmt.Println("Logs cleared successfully")
				return
			}

			f, err := os.Open(logs.LOG_FILENAME)
			if err != nil {
				fmt.Println("Error opening log file:", err)
				return
			}
			defer f.Close()

			finfo, err := f.Stat()
			if err != nil {
				fmt.Println("Error getting log file info:", err)
				return
			}

			content := make([]byte, finfo.Size())
			_, err = f.Read(content)
			if err != nil {
				fmt.Println("Error reading log file:", err)
				return
			}
			fmt.Println(string(content))
			return

		case "version":
			fmt.Println("gofast version:", VERSION)
			return
		}
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
