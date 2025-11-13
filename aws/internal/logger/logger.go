package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LogLevel representa el nivel de logging
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelError LogLevel = "ERROR"
)

// Logger interfaz para logging estructurado
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, err error, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
	WithRequestID(requestID string) Logger
}

// structuredLogger implementación de Logger con formato JSON
type structuredLogger struct {
	requestID string
	logLevel  LogLevel
}

// LogEntry representa una entrada de log en formato JSON
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	RequestID string                 `json:"request_id,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// New crea una nueva instancia de Logger
func New() Logger {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}
	return &structuredLogger{
		logLevel: LogLevel(logLevel),
	}
}

// WithRequestID crea un nuevo logger con el request ID especificado
func (l *structuredLogger) WithRequestID(requestID string) Logger {
	return &structuredLogger{
		requestID: requestID,
		logLevel:  l.logLevel,
	}
}

// Info registra un mensaje de nivel INFO
func (l *structuredLogger) Info(msg string, fields map[string]interface{}) {
	if l.shouldLog(LevelInfo) {
		l.log(LevelInfo, msg, "", fields)
	}
}

// Error registra un mensaje de nivel ERROR
func (l *structuredLogger) Error(msg string, err error, fields map[string]interface{}) {
	if l.shouldLog(LevelError) {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		l.log(LevelError, msg, errMsg, fields)
	}
}

// Debug registra un mensaje de nivel DEBUG
func (l *structuredLogger) Debug(msg string, fields map[string]interface{}) {
	if l.shouldLog(LevelDebug) {
		l.log(LevelDebug, msg, "", fields)
	}
}

// log escribe una entrada de log en formato JSON a stdout
func (l *structuredLogger) log(level LogLevel, msg string, errMsg string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     string(level),
		Message:   msg,
		RequestID: l.requestID,
		Error:     errMsg,
		Fields:    fields,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		// Fallback a log simple si falla la serialización JSON
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}

	fmt.Fprintln(os.Stdout, string(jsonBytes))
}

// shouldLog determina si un mensaje debe ser registrado según el nivel configurado
func (l *structuredLogger) shouldLog(level LogLevel) bool {
	levelPriority := map[LogLevel]int{
		LevelDebug: 0,
		LevelInfo:  1,
		LevelError: 2,
	}

	configuredPriority := levelPriority[l.logLevel]
	messagePriority := levelPriority[level]

	return messagePriority >= configuredPriority
}
