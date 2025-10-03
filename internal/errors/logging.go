package errors

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ErrorLogger provides structured error logging
type ErrorLogger struct {
	logger *log.Logger
	level  LogLevel
}

// NewErrorLogger creates a new error logger
func NewErrorLogger(level LogLevel) *ErrorLogger {
	return &ErrorLogger{
		logger: log.New(os.Stderr, "", 0),
		level:  level,
	}
}

// LogError logs an error with structured information
func (l *ErrorLogger) LogError(err error, context map[string]interface{}) {
	if err == nil {
		return
	}

	// Determine log level based on error type
	level := l.getLogLevel(err)
	if level < l.level {
		return
	}

	// Format the log message
	message := l.formatLogMessage(err, context, level)
	l.logger.Println(message)
}

// LogErrorf logs a formatted error message
func (l *ErrorLogger) LogErrorf(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	message := fmt.Sprintf("[%s] %s", level.String(), fmt.Sprintf(format, args...))
	l.logger.Println(message)
}

// LogGopherError logs a GopherError with full context
func (l *ErrorLogger) LogGopherError(err *GopherError, context map[string]interface{}) {
	if err == nil {
		return
	}

	level := l.getLogLevel(err)
	if level < l.level {
		return
	}

	// Merge context
	mergedContext := make(map[string]interface{})
	for k, v := range err.Context {
		mergedContext[k] = v
	}
	for k, v := range context {
		mergedContext[k] = v
	}

	message := l.formatGopherErrorLog(err, mergedContext, level)
	l.logger.Println(message)
}

// getLogLevel determines the appropriate log level for an error
func (l *ErrorLogger) getLogLevel(err error) LogLevel {
	if gopherErr, ok := err.(*GopherError); ok {
		switch gopherErr.Code {
		case ErrCodeInvalidVersion, ErrCodeInvalidArgument, ErrCodeInvalidFormat,
			ErrCodeMissingArgument, ErrCodeInvalidAliasName, ErrCodeReservedName,
			ErrCodeUnknownConfigOption, ErrCodeInvalidConfigValue:
			return LogLevelWarn // User errors are warnings
		case ErrCodeVersionNotInstalled, ErrCodeVersionAlreadyInstalled:
			return LogLevelInfo // These are informational
		case ErrCodeSystemGoNotAvailable, ErrCodeSymlinkFailed, ErrCodeEnvironmentSetupFailed,
			ErrCodeShellDetectionFailed, ErrCodePermissionDenied, ErrCodeDiskSpaceExhausted:
			return LogLevelError // System errors are errors
		case ErrCodeNetworkUnavailable, ErrCodeTimeoutExceeded, ErrCodeServerUnavailable,
			ErrCodeDownloadFailed:
			return LogLevelError // Network errors are errors
		case ErrCodeNotImplemented:
			return LogLevelWarn // Not implemented is a warning
		case ErrCodeOperationCancelled:
			return LogLevelInfo // Cancelled operations are informational
		default:
			return LogLevelError
		}
	}
	return LogLevelError
}

// formatLogMessage formats a log message with timestamp and context
func (l *ErrorLogger) formatLogMessage(err error, context map[string]interface{}, level LogLevel) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf("[%s] [%s] %s", timestamp, level.String(), err.Error())

	if len(context) > 0 {
		contextStr := l.formatContext(context)
		message += fmt.Sprintf(" | Context: %s", contextStr)
	}

	return message
}

// formatGopherErrorLog formats a GopherError for logging
func (l *ErrorLogger) formatGopherErrorLog(err *GopherError, context map[string]interface{}, level LogLevel) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, level.String(), err.Code, err.Message)

	if err.Details != "" {
		message += fmt.Sprintf(" | Details: %s", err.Details)
	}

	if len(context) > 0 {
		contextStr := l.formatContext(context)
		message += fmt.Sprintf(" | Context: %s", contextStr)
	}

	if err.File != "" {
		message += fmt.Sprintf(" | Location: %s:%d", err.File, err.Line)
	}

	if err.WrappedErr != nil {
		message += fmt.Sprintf(" | Wrapped: %v", err.WrappedErr)
	}

	return message
}

// formatContext formats context map as a string
func (l *ErrorLogger) formatContext(context map[string]interface{}) string {
	if len(context) == 0 {
		return ""
	}

	var parts []string
	for key, value := range context {
		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}
	return fmt.Sprintf("{%s}", fmt.Sprintf("%s", parts))
}

// SetLevel sets the minimum log level
func (l *ErrorLogger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel returns the current log level
func (l *ErrorLogger) GetLevel() LogLevel {
	return l.level
}

// Global logger instance
var (
	DefaultLogger = NewErrorLogger(LogLevelInfo)
)

// LogError is a convenience function that uses the default logger
func LogError(err error, context map[string]interface{}) {
	DefaultLogger.LogError(err, context)
}

// LogErrorf is a convenience function that uses the default logger
func LogErrorf(level LogLevel, format string, args ...interface{}) {
	DefaultLogger.LogErrorf(level, format, args...)
}

// LogGopherError is a convenience function that uses the default logger
func LogGopherError(err *GopherError, context map[string]interface{}) {
	DefaultLogger.LogGopherError(err, context)
}
