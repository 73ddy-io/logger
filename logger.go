package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log message.
type LogLevel int

// Available log levels.
//
// INFO is used for general informational messages.
// WARN is used for non-critical issues that might require attention.
// ERROR is used for errors and failures.
const (
	INFO LogLevel = iota
	WARN
	ERROR
)

var (
	logFile    *os.File
	logger     *log.Logger
	once       sync.Once
	timeFormat = "2006-01-02 15:04:05"
)

// InitLogger initializes the global logger with the given log file path.
//
// The function ensures that the directory for the log file exists,
// creates it if necessary, and then opens the log file for appending.
// If the logger is successfully initialized, subsequent logging
// functions (Info, Warn, Error) will write to this file.
func InitLogger(filename string) error {
	var err error
	fmt.Println("---------")

	dir := filepath.Dir(filename)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logFile, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("log error:", err.Error())
		return err
	}

	logger = log.New(logFile, "", 0)
	fmt.Println("---------")
	return nil
}

// Close closes the underlying log file if it is open.
//
// It should be called when the application is shutting down
// to ensure that all buffered log data is flushed and the
// file descriptor is released.
func Close() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

// Log writes a formatted log entry with the given level and message.
//
// Log automatically captures information about the caller (file name,
// line number, and function name) and prepends a timestamp and process ID.
// It is the low-level logging function that is wrapped by Info, Warn, and Error.
func Log(level LogLevel, message string) {
	if logger == nil {
		return
	}

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}

	shortFile := file
	if lastSlash := strings.LastIndex(file, "/"); lastSlash >= 0 {
		shortFile = file[lastSlash+1:]
	}

	funcName := runtime.FuncForPC(pc).Name()
	if lastDot := strings.LastIndex(funcName, "."); lastDot >= 0 {
		funcName = funcName[lastDot+1:]
	}

	levelStr := ""
	switch level {
	case INFO:
		levelStr = "INFO"
	case WARN:
		levelStr = "WARN"
	case ERROR:
		levelStr = "ERR"
	}

	pid := os.Getpid()

	logMsg := fmt.Sprintf("%s [%s] (%d)%s:%d %s - %s",
		time.Now().Format(timeFormat),
		levelStr,
		pid,
		shortFile,
		line,
		funcName,
		message,
	)

	logger.Println(logMsg)
}

// Info logs an informational message using printf-style formatting.
//
// The format string and arguments are passed to fmt.Sprintf
// and the resulting string is logged with INFO level.
func Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Log(INFO, message)
}

// Warn logs a warning message using printf-style formatting.
//
// This should be used for situations that are not fatal but may
// require attention or indicate a potential problem.
func Warn(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Log(WARN, message)
}

// Error logs an error message using printf-style formatting.
//
// Use this for error conditions and failures that should be visible
// in application logs.
func Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Log(ERROR, message)
}
