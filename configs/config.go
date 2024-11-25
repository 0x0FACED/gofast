package configs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	NumWorkers     int
	ResultChanSize int
	ErrorsChanSize int
	LogFilename    string
}

type cfg struct {
	NumWorkers     string
	ResultChanSize string
	ErrorsChanSize string
	LogFilename    string
}

func Load() (*AppConfig, error) {
	configFile := "app.conf"

	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	temp := cfg{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "NUM_WORKERS":
			temp.NumWorkers = value
		case "RESULT_CHAN_SIZE":
			temp.ResultChanSize = value
		case "ERRORS_CHAN_SIZE":
			temp.ErrorsChanSize = value
		case "LOG_FILENAME":
			temp.LogFilename = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	config, err := parse(temp)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func parse(temp cfg) (*AppConfig, error) {
	workers, err := strconv.Atoi(temp.NumWorkers)
	if err != nil {
		return nil, fmt.Errorf("invalid value for NUM_WORKERS: %v", err)
	}

	resultChanSize, err := strconv.Atoi(temp.ResultChanSize)
	if err != nil {
		return nil, fmt.Errorf("invalid value for RESULT_CHAN_SIZE: %v", err)
	}

	errorsChanSize, err := strconv.Atoi(temp.ErrorsChanSize)
	if err != nil {
		return nil, fmt.Errorf("invalid value for ERRORS_CHAN_SIZE: %v", err)
	}

	return &AppConfig{
		NumWorkers:     workers,
		ResultChanSize: resultChanSize,
		ErrorsChanSize: errorsChanSize,
		LogFilename:    temp.LogFilename,
	}, nil
}
