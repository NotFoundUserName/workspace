package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type LogLevel string

const (
	INFO  LogLevel = "info"
	ERROR LogLevel = "error"
	WARN  LogLevel = "warn"
	DEBUG LogLevel = "debug"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"`
	Error     string    `json:"error,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", 0)
}

func logWithLevel(level LogLevel, message string, err error, context map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}

	if err != nil {
		entry.Error = err.Error()
	}

	if context != nil {
		entry.Context = context
	}

	data, _ := json.Marshal(entry)
	logger.Println(string(data))
}

func Info(message string, context ...map[string]interface{}) {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	logWithLevel(INFO, message, nil, ctx)
}

func Error(message string, err error, context ...map[string]interface{}) {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	logWithLevel(ERROR, message, err, ctx)
}

func Warn(message string, context ...map[string]interface{}) {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	logWithLevel(WARN, message, nil, ctx)
}

func Debug(message string, context ...map[string]interface{}) {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	logWithLevel(DEBUG, message, nil, ctx)
}
