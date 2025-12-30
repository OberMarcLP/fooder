package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	debugMode bool
	Logger    zerolog.Logger
)

// Context key for request ID
type contextKey string

const RequestIDKey contextKey = "request_id"

func init() {
	// Check if debug mode is enabled
	debugMode = os.Getenv("DEBUG") == "true"

	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339

	// Set log level
	if debugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Use pretty console output in development
	if os.Getenv("LOG_FORMAT") != "json" {
		Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		})
	} else {
		Logger = log.Logger
	}
}

// Debug logs a debug message with optional fields
func Debug(format string, v ...interface{}) {
	Logger.Debug().Msgf(format, v...)
}

// DebugWithFields logs a debug message with structured fields
func DebugWithFields(msg string, fields map[string]interface{}) {
	event := Logger.Debug()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	Logger.Info().Msgf(format, v...)
}

// InfoWithFields logs an info message with structured fields
func InfoWithFields(msg string, fields map[string]interface{}) {
	event := Logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Warn logs a warning message
func Warn(format string, v ...interface{}) {
	Logger.Warn().Msgf(format, v...)
}

// WarnWithFields logs a warning message with structured fields
func WarnWithFields(msg string, fields map[string]interface{}) {
	event := Logger.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	Logger.Error().Msgf(format, v...)
}

// ErrorWithFields logs an error message with structured fields
func ErrorWithFields(msg string, fields map[string]interface{}) {
	event := Logger.Error()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Fatal logs a fatal message and exits
func Fatal(format string, v ...interface{}) {
	Logger.Fatal().Msgf(format, v...)
}

// IsDebugMode returns whether debug mode is enabled
func IsDebugMode() bool {
	return debugMode
}

// WithRequestID creates a context-aware logger with request ID
func WithRequestID(ctx context.Context) *zerolog.Logger {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		logger := Logger.With().Str("request_id", requestID).Logger()
		return &logger
	}
	return &Logger
}

// LogRequest logs an HTTP request with structured data
func LogRequest(method, path, requestID, ip string, duration time.Duration, status int, bytes int64) {
	event := Logger.Info()

	if requestID != "" {
		event = event.Str("request_id", requestID)
	}

	event.
		Str("method", method).
		Str("path", path).
		Str("ip", ip).
		Dur("duration", duration).
		Int("status", status).
		Int64("bytes", bytes).
		Msg("HTTP request completed")
}

// LogError logs an error with additional context
func LogError(err error, msg string, fields map[string]interface{}) {
	event := Logger.Error().Err(err)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// LogDatabaseQuery logs a database query with performance metrics
func LogDatabaseQuery(query string, duration time.Duration, rows int) {
	if debugMode {
		Logger.Debug().
			Str("query", query).
			Dur("duration", duration).
			Int("rows", rows).
			Msg("Database query executed")
	}
}
