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

// LogLevel represents the severity of the log message.
type LogLevel int

// Constants defining available log levels.
const (
	INFO LogLevel = iota
	WARN
	ERROR
)

// Internal package variables.
var (
	logFile    *os.File
	logger     *log.Logger
	once       sync.Once
	timeFormat = "2006-01-02 15:04:05"
)

// InitLogger initializes the logger instance with the specified file path.
//
// @brief Sets up the log directory and file, creating them if necessary.
// @param filename The path to the log file (e.g., "logs/app.log").
// @return error Returns an error if directory creation or file opening fails.
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

// Close gracefully closes the active log file.
//
// @brief Closes the underlying file handle if it is open.
// @return error Returns an error if the file close operation fails.
func Close() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

// Log is the generic logging function that handles formatting and writing.
//
// @brief Captures runtime caller info and writes a formatted log entry.
// @param level The severity level of the log (INFO, WARN, ERROR).
// @param message The actual log message string.
func Log(level LogLevel, message string) {
	if logger == nil {
		return
	}

	// Get information about the calling function (stack depth 2)
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}

	// Extract only the filename from the full path
	shortFile := file
	if lastSlash := strings.LastIndex(file, "/"); lastSlash >= 0 {
		shortFile = file[lastSlash+1:]
	}

	// Get the function name
	funcName := runtime.FuncForPC(pc).Name()
	if lastDot := strings.LastIndex(funcName, "."); lastDot >= 0 {
		funcName = funcName[lastDot+1:]
	}

	// Determine the string representation of the log level
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

	// Format the final log string
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

// Info logs an informational message.
//
// @brief Wrapper for the generic Log function with INFO level.
// @param format The format string (printf style).
// @param args The arguments for the format string.
func Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Log(INFO, message)
}

// Warn logs a warning message.
//
// @brief Wrapper for the generic Log function with WARN level.
// @param format The format string (printf style).
// @param args The arguments for the format string.
func Warn(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Log(WARN, message)
}

// Error logs an error message.
//
// @brief Wrapper for the generic Log function with ERROR level.
// @param format The format string (printf style).
// @param args The arguments for the format string.
func Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Log(ERROR, message)
}
