# gofast

The purpose of this utility is to overtake the `find` utility.


## Flags

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

## Run

```sh
gofast -p / -t [file | dir] -n "filenameOrDir" -m [exact | pattern]
```

## Alternative methods

### Clone, build and run

You can clone repo, build executable file and run with the following steps:

1. Clone repo using HTTPS:

```sh
git clone https://github.com/0x0FACED/gofast.git
```

2. Build executable

```sh
go build -o gofast cmd/gofast/main.go
```

3. Run

```sh
./gofast -p / -t file -n "gofast" -m exact
```

### Check version

```sh
gofast version
```

### Read logs

```sh
gofast logs
```

### Clear logs

```sh
gofast logs clear
```

### Work in progress