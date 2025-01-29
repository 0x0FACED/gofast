package logs

import (
	"os"
	"time"
)

const (
	LOG_FILENAME = "logs.log"
)

type Logger struct {
	logFile *os.File
}

func New() (*Logger, error) {
	file, err := os.OpenFile(LOG_FILENAME, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		logFile: file,
	}, nil
}

func (l *Logger) Close() {
	l.logFile.Close()
}

func (l *Logger) Info(msg string) {
	l.logFile.Write([]byte(time.Now().Format("[2006.01.02 | 15:04:06]") + "\t[INFO]\t" + msg + "\n"))
}

func (l *Logger) Error(err error) {
	l.logFile.Write([]byte(time.Now().Format("[2006.01.02 | 15:04:06]") + "\t[ERROR]\t" + err.Error() + "\n"))
}

func (l *Logger) Debug(msg string) {
	l.logFile.Write([]byte(time.Now().Format("[2006.01.02 | 15:04:06]") + "\t[DEBUG]\t" + msg + "\n"))
}
