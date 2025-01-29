# gofast

The purpose of this utility is to overtake the `find` utility.

## Run

Run with command:

```sh
sudo go run path/to/main.go -p / -t file -n "filename" -m "exact"
```

Use `sudo` to make the search full-fledged.

```go
flag.StringVar(&startPath, "p", "/", "Starting point for the search (required parameter)")
flag.StringVar(&fileType, "t", "file", "Type of search: 'file' or 'dir' (required parameter)")
flag.StringVar(&name, "n", "", "Name of the file or directory in quotes (required parameter)")
flag.StringVar(&method, "m", "like", "Search method, e.g., 'exact' or 'pattern'")
```

## Install

You can install the application using the command:
```sh
go install github.com/0x0FACED/gofast/cmd/gofast@latest
```

## Usage after install

### Run

```sh
gofast -p / -t [file | dir] -n "filenameOrDir" -m [exact | pattern]
```

### Read logs

```sh
gofast logs
```

### Work in progress