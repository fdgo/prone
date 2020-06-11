package loggers

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

type LOG_LEVEL int

const (
	LOG_LEVEL_UNKOWN LOG_LEVEL = iota
	LOG_LEVEL_ERROR            // 错误级别
	LOG_LEVEL_WARN             // 警告级别
	LOG_LEVEL_DEBUG            // 调试级别
	LOG_LEVEL_INFO             // 信息级别
)

var (
	Error Logger
	Info  Logger
	Debug Logger
	Warn  Logger
)

func init() {
	Error = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	Info = log.New(os.Stderr, "[Info] ", log.LstdFlags|log.Lshortfile)
	Warn = log.New(os.Stderr, "[Warn] ", log.LstdFlags|log.Lshortfile)
	Debug = log.New(os.Stderr, "[Debug] ", log.LstdFlags|log.Lshortfile)
}

func InitLogLevel(level string) {
	switch strings.ToUpper(level) {
	case "ERROR":
		Warn = new(NullLogger)
		fallthrough
	case "WARN":
		Info = new(NullLogger)
		fallthrough
	case "INFO":
		Debug = new(NullLogger)
	}
}

func SetOutput(w io.Writer) {
	Error.SetOutput(w)
	Info.SetOutput(w)
	Debug.SetOutput(w)
	Warn.SetOutput(w)
}

func CaptureLog(f func()) string {
	var buf bytes.Buffer
	SetOutput(&buf)
	f()
	SetOutput(os.Stderr)
	return buf.String()
}

// Logger 接口
type Logger interface {
	SetOutput(w io.Writer)
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

// NullLogger 空日志
type NullLogger int

// SetOutput sets the output destination for the logger.
func (l *NullLogger) SetOutput(w io.Writer) {}

// Printf calls l.Output to print to the logger.
func (l *NullLogger) Printf(format string, v ...interface{}) {}

// Print calls l.Output to print to the logger.
func (l *NullLogger) Print(v ...interface{}) {}

// Println calls l.Output to print to the logger.
func (l *NullLogger) Println(v ...interface{}) {}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *NullLogger) Fatal(v ...interface{}) {}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *NullLogger) Fatalf(format string, v ...interface{}) {}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *NullLogger) Fatalln(v ...interface{}) {}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *NullLogger) Panic(v ...interface{}) {}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *NullLogger) Panicf(format string, v ...interface{}) {}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *NullLogger) Panicln(v ...interface{}) {}
