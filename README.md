# gofast

The purpose of this utility is to overtake the find utility.

## Run

Run with command:

```sh
sudo go run path/to/main.go -p / -t file -n "filename" -m like -workers 100
```

Use `sudo` to make the search full-fledged.

```go
flag.StringVar(&startPath, "p", "/", "Starting point for the search (required parameter)")
flag.StringVar(&fileType, "t", "file", "Type of search: 'file' or 'dir' (required parameter)")
flag.StringVar(&name, "n", "", "Name of the file or directory in quotes (required parameter)")
flag.StringVar(&method, "m", "like", "Search method, e.g., 'like' or 'pattern'")
flag.IntVar(&workers, "workers", 10, "Number of goroutines for parallel search (default is 10)")
```

It's running slower at the moment. But there are some plans for optimization. The main one is to optimize the work of workers.

### Work in progress