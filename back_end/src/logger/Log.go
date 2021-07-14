package logger

import (
	"fmt"
	"github.com/kr/pretty"
	"io"
	"log"
	"os"
	"rabbit_gather/util"
)

const (
	DEBUG uint8 = 1 << iota
	WARNING
	ERROR
)
const (
	MUTE = uint8(0)
	ALL  = ^MUTE
)

var LogLevelMask = ALL

func init() {
	type Config struct {
		MinLogLevel int16 `json:"log_level"`
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/log.config.json")
	if err != nil {
		panic(err.Error())
	}
	if config.MinLogLevel != -1 {
		LogLevelMask = uint8(config.MinLogLevel)
	}

}

type LoggerWrapper struct {
	DEBUG      *Logger
	WARNING    *Logger
	ERROR      *Logger
	tempLogger *Logger
}

type Logger struct {
	log.Logger
}

func (l *Logger) PrettyPrintln(v ...interface{}) {
	l.Println(pretty.Sprint(v...))
}

const TempLogOpen = true

func (l *LoggerWrapper) TempLog() *Logger {
	if l.tempLogger == nil {
		if !TempLogOpen {
			l.tempLogger = &Logger{*log.New(io.Discard, "TEMP_LOG: ", log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix)}
		} else {
			l.tempLogger = &Logger{*log.New(os.Stdout, "TEMP_LOG: ", log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix)}
		}
	}
	return l.tempLogger
	//if len(things) ==1{
	//	l.tempLogger.Println(things[0])
	//	return
	//}else{
	//	l.tempLogger.Printf(fmt.Sprintf("%s\n",fmt.Sprint(things[0])),things[1:]...)
	//	return
	//}

}

func NewLoggerWrapper(prefix string) *LoggerWrapper {
	fmt.Printf("Cteate logger: %s\n", prefix)
	return &LoggerWrapper{
		ERROR:   CreateErrorLogger(prefix),
		WARNING: CreateWaringLogger(prefix),
		DEBUG:   CreateDebugLogger(prefix),
	}
}

func CreateErrorLogger(prefix string) *Logger {
	outputFile, err := os.OpenFile(fmt.Sprintf("../log/error_rabbit_gather_%s.log", prefix), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	var writer io.Writer
	if ERROR&LogLevelMask == 0 {
		writer = io.Discard
	} else {
		writer = io.MultiWriter(outputFile, os.Stdout)
	}
	return &Logger{*log.New(writer, fmt.Sprintf("%s - ERROR: ", prefix), log.Lmicroseconds|log.Ldate|log.Llongfile|log.Lmsgprefix)}
}
func CreateWaringLogger(prefix string) *Logger {
	outputFile, err := os.OpenFile(fmt.Sprintf("../log/warning_rabbit_gather_%s.log", prefix), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	var writer io.Writer
	if WARNING&LogLevelMask == 0 {
		writer = io.Discard
	} else {
		writer = io.MultiWriter(outputFile, os.Stdout)
	}
	return &Logger{*log.New(writer, fmt.Sprintf("%s - WARNING: ", prefix), log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix)}

}
func CreateDebugLogger(prefix string) *Logger {
	var writer io.Writer
	if DEBUG&LogLevelMask == 0 {
		writer = io.Discard
	} else {
		writer = os.Stdout
	}
	return &Logger{*log.New(writer, fmt.Sprintf("%s - DEBUG: ", prefix), log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)}
}
