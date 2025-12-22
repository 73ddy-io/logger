# logger

Simple structured logging library for Go applications. Supports INFO, WARN, ERROR levels, file output, timestamps, caller info, and PID.

## Installation

```bash
go get github.com/73ddy-io/logger
```

## Usage

```go
package main

import (
    "github.com/73ddy-io/logger"
)

func main() {
    logger.InitLogger("app.log")
    defer logger.Close()
    
    logger.Info("Application started")
    logger.Warn("Warning: low memory")
    logger.Error("Critical error: %s", "database unavailable")
}
```

## Log Format

```
2025-01-02 15:04:05 [INFO] 1234 main.go:12 main.main - Application started
```

## Functions

- `InitLogger(filename string) error` — initializes logger with file
- `Close() error` — closes log file
- `Info(format string, args ...interface{})`
- `Warn(format string, args ...interface{})`
- `Error(format string, args ...interface{})`

## License

MIT