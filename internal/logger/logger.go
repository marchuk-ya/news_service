package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents log level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	RequestID   string                 `json:"request_id,omitempty"`
	Method      string                 `json:"method,omitempty"`
	Path        string                 `json:"path,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
	Environment string                 `json:"environment,omitempty"`
}

// Logger provides structured logging functionality
type Logger struct {
	level       Level
	environment string
}

// New creates a new logger instance
func New(level string) *Logger {
	var logLevel Level
	switch level {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
	case "warn":
		logLevel = WARN
	case "error":
		logLevel = ERROR
	default:
		logLevel = INFO
	}

	return &Logger{
		level:       logLevel,
		environment: os.Getenv("ENVIRONMENT"),
	}
}

// WithRequest creates a logger with request context
func (l *Logger) WithRequest(method, path, userAgent string) *Logger {
	return &Logger{
		level:       l.level,
		environment: l.environment,
	}
}

// log writes a log entry
func (l *Logger) log(level Level, message string, fields ...interface{}) {
	if level < l.level {
		return
	}

	entry := LogEntry{
		Timestamp:   time.Now(),
		Level:       level.String(),
		Message:     message,
		Environment: l.environment,
		Fields:      make(map[string]interface{}),
	}

	// Parse fields
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			if key, ok := fields[i].(string); ok {
				entry.Fields[key] = fields[i+1]
			}
		}
	}

	// Marshal to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}

	fmt.Println(string(data))
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...interface{}) {
	l.log(DEBUG, message, fields...)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...interface{}) {
	l.log(INFO, message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...interface{}) {
	l.log(WARN, message, fields...)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...interface{}) {
	l.log(ERROR, message, fields...)
}

// WithContext creates a logger with context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract request ID from context if available
	if _, ok := ctx.Value("request_id").(string); ok {
		return &Logger{
			level:       l.level,
			environment: l.environment,
		}
	}

	return &Logger{
		level:       l.level,
		environment: l.environment,
	}
}
