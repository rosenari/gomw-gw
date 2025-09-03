package logger

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type Fields map[string]interface{}

var (
	instance *JSONLogger
	once     sync.Once
)

type JSONLogger struct {
	logger *log.Logger
}

func GetLogger() *JSONLogger {
	once.Do(func() {
		instance = &JSONLogger{
			logger: log.New(os.Stdout, "", 0),
		}
	})
	return instance
}

func Debug(message string, fields Fields) {
	GetLogger().writeLog("debug", message, fields)
}

func Info(message string, fields Fields) {
	GetLogger().writeLog("info", message, fields)
}

func Warn(message string, fields Fields) {
	GetLogger().writeLog("warn", message, fields)
}

func Error(message string, fields Fields) {
	GetLogger().writeLog("error", message, fields)
}

func Fatal(message string, fields Fields) {
	GetLogger().writeLog("fatal", message, fields)
	os.Exit(1)
}

func (l *JSONLogger) writeLog(level, message string, fields Fields) {
	logEntry := map[string]interface{}{
		"level":     level,
		"message":   message,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	for key, value := range fields {
		logEntry[key] = value
	}

	if err := json.NewEncoder(os.Stdout).Encode(logEntry); err != nil {
		log.Printf("[%s] %s: %v", level, message, fields)
	}
} 