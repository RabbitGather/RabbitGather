package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	ERROR   = "ERROR"
	WARNING = "WARNING"
	DEBUG   = "DEBUG"
)

type LoggerWrapper struct {
	DEBUG      *log.Logger
	WARNING    *log.Logger
	ERROR      *log.Logger
	tempLogger *log.Logger
}

func (l *LoggerWrapper) tempLog(things ...interface{}) {
	if l.tempLogger == nil {
		l.tempLogger = log.New(os.Stdout, "TEMP_LOG: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
	}
	l.tempLogger.Println(things...)
}
func NewLogger(prefix string) LoggerWrapper {
	return LoggerWrapper{
		ERROR:   CreateLogger(prefix, ERROR),
		WARNING: CreateLogger(prefix, WARNING),
		DEBUG:   CreateLogger(prefix, DEBUG),
	}
}
func CreateLogger(prefix, loglevel string) *log.Logger {
	var logger *log.Logger
	//var outputFile *os.File
	switch loglevel {
	case ERROR:
		outputFile, err := os.OpenFile(fmt.Sprintf("log/error_rabbit_gather_%s.log", prefix), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err.Error())
		}
		writer := io.MultiWriter(outputFile, os.Stdout)
		logger = log.New(writer, fmt.Sprintf("ERROR - %s: ", prefix), log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
	case WARNING:
		outputFile, err := os.OpenFile(fmt.Sprintf("log/warning_rabbit_gather_%s.log", prefix), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err.Error())
		}
		writer := io.MultiWriter(outputFile, os.Stdout)
		logger = log.New(writer, fmt.Sprintf("WARNING - %s: ", prefix), log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix)
	case DEBUG:
		logger = log.New(os.Stdout, fmt.Sprintf("DEBUG - %s: ", prefix), log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix)
	default:
		panic("loglevel wrong")
	}

	return logger
}
