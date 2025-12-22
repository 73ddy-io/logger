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

type LogLevel int

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

// InitLogger инициализирует логгер с указанным файлом
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

// Close закрывает файл лога
func Close() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

// Log универсальная функция логирования
func Log(level LogLevel, message string) {
	if logger == nil {
		return
	}

	// Получаем информацию о вызывающей функции
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}

	// Извлекаем только имя файла из полного пути
	shortFile := file
	if lastSlash := strings.LastIndex(file, "/"); lastSlash >= 0 {
		shortFile = file[lastSlash+1:]
	}

	// Получаем имя функции
	funcName := runtime.FuncForPC(pc).Name()
	if lastDot := strings.LastIndex(funcName, "."); lastDot >= 0 {
		funcName = funcName[lastDot+1:]
	}

	// Формируем уровень логирования
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
	// Формируем итоговую строку лога
	logMsg := fmt.Sprintf("%s [%s] (%d)%s:%d %s - %s",
		time.Now().Format(timeFormat),
		levelStr,
		shortFile,
		pid,
		line,
		funcName,
		message,
	)

	logger.Println(logMsg)
}

// Info логирование информационного сообщения
func Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Log(INFO, message)
}

// Warn логирование предупреждения
func Warn(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Log(WARN, message)
}

// Error логирование ошибки
func Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Log(ERROR, message)
}
