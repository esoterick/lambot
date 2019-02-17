package log

import (
	"fmt"
	"log"
	"os"
	"time"

	"path/filepath"
	"runtime"

	"github.com/natefinch/lumberjack"
)

var (
	// ErrorLog logger with custom formatting, written to stdOut
	ErrorLog *log.Logger
	// FatalLog logger with custom formatting, written to stdOut
	FatalLog *log.Logger
	// InfoLog logger with custom formatting, written to stdOut
	InfoLog *log.Logger
	// WarningLog logger with custom formatting, written to stdOut
	WarningLog *log.Logger
	// Color
	color bool
)

// AppLogLocation is our log file name
const AppLogLocation = "lambot.log"

// ColorRed ANSI Color Red
const ColorRed = "\x1b[31;1m"

// ColorYellow ANSI Color Yellow
const ColorYellow = "\x1b[33;1m"

// ColorWhite ANSI Color White
const ColorWhite = "\x1b[37;1m"

// ColorClear ANSI Clear
const ColorClear = "\x1b[0m"

func init() {
	applicationLogFile := &lumberjack.Logger{
		Filename:   AppLogLocation,
		MaxAge:     1,  // days
		MaxBackups: 14, // max days of retention
	}

	_, ok := os.LookupEnv("DEBUG")
	if ok {
		InfoLog = log.New(os.Stdout, "", 0)
		WarningLog = log.New(os.Stdout, "", 0)
		ErrorLog = log.New(os.Stdout, "", 0)
		FatalLog = log.New(os.Stdout, "", 0)

		_, disable := os.LookupEnv("DISABLE_COLOR")
		if disable {
			color = false
		} else {
			color = true
		}
	} else {
		InfoLog = log.New(applicationLogFile, "", 0)
		WarningLog = log.New(applicationLogFile, "", 0)
		ErrorLog = log.New(applicationLogFile, "", 0)
		FatalLog = log.New(applicationLogFile, "", 0)
		color = false
	}
}

func getCaller() string {
	pc, fn, line, _ := runtime.Caller(2)
	_, file := filepath.Split(fn)
	_, fc := filepath.Split(runtime.FuncForPC(pc).Name())
	return fmt.Sprintf("%s [%s:%d] ", fc, file, line)
}

func colorMe(s string) string {
	pre := ""
	post := ""
	if color {
		switch s {
		case "[ERROR]", "[FATAL]":
			pre = ColorRed
		case "[WARN]":
			pre = ColorYellow
		case "[INFO]":
			pre = ColorWhite
		}
		post = ColorClear
	}
	return fmt.Sprintf("%s%s%s", pre, s, post)
}

// Error logs error messages with custom format using Print
func Error(m ...interface{}) {
	prefix := fmt.Sprintf("%-28s %-7s %s", time.Now().Format("01/02/2006 15:04:05.999 MST"), colorMe("[ERROR]"), getCaller())
	ErrorLog.Print(append([]interface{}{prefix}, m...)...)
}

// Errorf logs error messages with custom format using Printf
func Errorf(m string, v ...interface{}) {
	prefix := fmt.Sprintf("%-28s %-7s %s", time.Now().Format("01/02/2006 15:04:05.999 MST"), colorMe("[ERROR]"), getCaller())
	ErrorLog.Printf(prefix+m, v...)
}

// Fatal logs error messages with custom format using Print
func Fatal(m ...interface{}) {
	prefix := fmt.Sprintf("%-28s %-7s %s", time.Now().Format("01/02/2006 15:04:05.999 MST"), colorMe("[FATAL]"), getCaller())
	FatalLog.Print(append([]interface{}{prefix}, m...)...)
	os.Exit(1)
}

// Fatalf logs error messages with custom format using Printf
func Fatalf(m string, v ...interface{}) {
	prefix := fmt.Sprintf("%-28s %-7s %s", time.Now().Format("01/02/2006 15:04:05.999 MST"), colorMe("[FATAL]"), getCaller())
	FatalLog.Printf(prefix+m, v...)
	os.Exit(1)
}

// Info logs info messages with custom format using Print
func Info(m ...interface{}) {
	prefix := fmt.Sprintf("%-28s %-7s %s", time.Now().Format("01/02/2006 15:04:05.999 MST"), colorMe("[INFO]"), getCaller())
	InfoLog.Print(append([]interface{}{prefix}, m...)...)
}

// Infof logs info messages with custom format using Printf
func Infof(m string, v ...interface{}) {
	prefix := fmt.Sprintf("%-28s %-7s %s", time.Now().Format("01/02/2006 15:04:05.999 MST"), colorMe("[INFO]"), getCaller())
	InfoLog.Printf(prefix+m, v...)
}

// Warn logs warning messages with custom format using Print
func Warn(m ...interface{}) {
	prefix := fmt.Sprintf("%-28s %-7s %s", time.Now().Format("01/02/2006 15:04:05.999 MST"), colorMe("[WARN]"), getCaller())
	WarningLog.Print(append([]interface{}{prefix}, m...)...)
}

// Warnf logs warning messages with custom format using Printf
func Warnf(m string, v ...interface{}) {
	prefix := fmt.Sprintf("%-28s %-7s %s", time.Now().Format("01/02/2006 15:04:05.999 MST"), colorMe("[WARN]"), getCaller())
	WarningLog.Printf(prefix+m, v...)
}
