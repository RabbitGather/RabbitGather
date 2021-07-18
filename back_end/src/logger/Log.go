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

//const ResetColor =  "\033[0m"
//const RedColor =  "\033[97;41m"

var LogLevelMask = ALL
var DebugColor util.ColorCode
var WaringColor util.ColorCode
var ErrorColor util.ColorCode
var TempColor util.ColorCode
var TempLogOpen bool

func init() {
	type Config struct {
		MinLogLevel int16 `json:"log_level"`
		DebugColor  int   `json:"debug_color"`
		WaringColor int   `json:"waring_color"`
		ErrorColor  int   `json:"error_color"`
		TempColor   int   `json:"temp_color"`
		TempLogOpen bool  `json:"temp_log_open"`
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/log.config.json")
	if err != nil {
		panic(err.Error())
	}
	if config.MinLogLevel < -1 {
		panic("The log level must >= -1")
	}
	if config.MinLogLevel >= 0 {
		LogLevelMask = uint8(config.MinLogLevel)
	}
	DebugColor = util.ColorCode(config.DebugColor)
	WaringColor = util.ColorCode(config.WaringColor)
	ErrorColor = util.ColorCode(config.ErrorColor)
	TempColor = util.ColorCode(config.TempColor)
	TempLogOpen = config.TempLogOpen
}

type LoggerWrapper struct {
	DEBUG      *Logger
	WARNING    *Logger
	ERROR      *Logger
	tempLogger *Logger
}

type Logger struct {
	log.Logger
	Color util.ColorCode
}

func (l *Logger) PrettyPrintln(v ...interface{}) {
	l.Println(pretty.Sprint(v...))
}

func (l *Logger) SetColor(color util.ColorCode) {
	l.Color = color
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Output(2, fmt.Sprintf(util.ColorSting(format, l.Color), v...))
	//l.Logger.Printf(util.ColorSting(format,l.Color),v...)
}
func (l *Logger) Print(v ...interface{}) {
	l.Logger.Output(2, util.ColorSting(fmt.Sprint(v...), l.Color))
	//l.Logger.Print(util.ColorSting(fmt.Sprint(v...),l.Color))
}
func (l *Logger) Println(v ...interface{}) {
	l.Output(2, util.ColorSting(pretty.Sprint(v...), l.Color))
	//l.Logger.Println(util.ColorSting(fmt.Sprint(v...),l.Color))
}

func (l *LoggerWrapper) TempLog() *Logger {
	if l.tempLogger == nil {
		var writer io.Writer
		if TempLogOpen {
			writer = os.Stdout
		} else {
			writer = io.Discard
		}
		l.tempLogger = &Logger{*log.New(writer, util.ColorSting("TEMP_LOG: ", TempColor), log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix), TempColor}
	}
	return l.tempLogger
}

func NewLoggerWrapper(prefix string) *LoggerWrapper {
	if LogLevelMask != MUTE {
		fmt.Printf("Cteate logger: %s\n", prefix)
	}
	//fmt.Println("ERER")
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
	return &Logger{*log.New(writer, fmt.Sprintf(util.ColorSting("ERROR %s: ", ErrorColor), prefix), log.Lmicroseconds|log.Ldate|log.Llongfile|log.Lmsgprefix), ErrorColor}
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
	return &Logger{*log.New(writer, fmt.Sprintf(util.ColorSting("WARNING %s: ", WaringColor), prefix), log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix), WaringColor}

}
func CreateDebugLogger(prefix string) *Logger {
	var writer io.Writer
	if DEBUG&LogLevelMask == 0 {
		writer = io.Discard
	} else {
		writer = os.Stdout
	}
	return &Logger{*log.New(writer, fmt.Sprintf(util.ColorSting("DEBUG %s: ", DebugColor), prefix), log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix), DebugColor}
}
