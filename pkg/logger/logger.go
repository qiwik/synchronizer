package logger

import (
	"errors"
	"log"
	"os"
)

//Logger contains information about levels of logging
type Logger struct {
	Info *log.Logger
	Errs *log.Logger
	Fatal *log.Logger
}

// LogFileInit returns a txt log file
func LogFileInit() (*os.File, error) {
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, errors.New("error: can't work with log.txt")
	}
	return logFile, nil
}

// LogInit creates new loggers by levels
func LogInit(file *os.File) *Logger{
	var newLog Logger
	newLog.Info = log.New(file, "INFO\t", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	newLog.Errs = log.New(file, "ERROR\t", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	newLog.Fatal = log.New(file, "Fatal\t", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	return &newLog
}
