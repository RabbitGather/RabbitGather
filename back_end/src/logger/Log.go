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
	// Debug logger should be used to log a message that will be useful for develop.
	DEBUG *Logger
	// Waring logger should be used to log an error that will not break business logic.
	WARNING *Logger
	// error logger should be used to log an error that will break business logic.
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

// Printf will color the text and act like log.Printf()
func (l *Logger) Printf(format string, v ...interface{}) {
	_ = l.Logger.Output(2, fmt.Sprintf(util.ColorSting(format, l.Color), v...))
}

// Print will color the text and act like log.Print()
func (l *Logger) Print(v ...interface{}) {
	_ = l.Logger.Output(2, util.ColorSting(fmt.Sprint(v...), l.Color))
}

// Println will color the text and act like log.Println()
func (l *Logger) Println(v ...interface{}) {
	_ = l.Output(2, util.ColorSting(fmt.Sprint(v...), l.Color))
}

// The TempLog should only use for debug, it will be close if the TempLogOpen parameter is false
// Se the settings in config/log.config.json
func (l *LoggerWrapper) TempLog() *Logger {
	if l.tempLogger == nil {
		if !TempLogOpen {
			return CreateMuteLogger()
		}
		l.tempLogger = &Logger{*log.New(os.Stdout, util.ColorSting("TEMP_LOG: ", TempColor), log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix), TempColor}
	}
	return l.tempLogger
}

// NewLoggerWrapper Create a new LoggerWrapper with given prefix,
// The prefix will be print before all log rows
func NewLoggerWrapper(prefix string) *LoggerWrapper {
	if LogLevelMask == MUTE {
		return &LoggerWrapper{
			ERROR:   CreateMuteLogger(),
			WARNING: CreateMuteLogger(),
			DEBUG:   CreateMuteLogger(),
		}
	}
	fmt.Printf("Cteate logger: %s\n", prefix)
	return NewMuteLoggerWrapper()
}

// NewMuteLoggerWrapper create a mute logger that will do nothing when use
func NewMuteLoggerWrapper() *LoggerWrapper {
	return &LoggerWrapper{
		ERROR:   CreateMuteLogger(),
		WARNING: CreateMuteLogger(),
		DEBUG:   CreateMuteLogger(),
	}

}

// CreateMuteLogger create a Mute Logger, the mute logger will do nothing when used.
func CreateMuteLogger() *Logger {
	return &Logger{*log.New(io.Discard, "", log.LstdFlags), util.FgBlack}
}

// CreateErrorLogger create an Error Logger.
// error logger should be used to log an error that will break business logic.
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

// CreateWaringLogger create a Waring Logger.
// Waring logger should be used to log an error that will not break business logic.
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

// CreateDebugLogger create a Waring Logger.
// Debug logger should be used to log a message that will be useful for develop.
func CreateDebugLogger(prefix string) *Logger {
	var writer io.Writer
	if DEBUG&LogLevelMask == 0 {
		writer = io.Discard
	} else {
		writer = os.Stdout
	}
	return &Logger{*log.New(writer, fmt.Sprintf(util.ColorSting("DEBUG %s: ", DebugColor), prefix), log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix), DebugColor}
}
