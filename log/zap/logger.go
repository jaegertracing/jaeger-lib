package zap

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger is an adapter from zap Logger to jaeger-lib Logger.
type Logger struct {
	logger zap.Logger
}

// NewLogger creates a new Logger.
func NewLogger(logger zap.Logger) *Logger {
	return &Logger{logger: logger}
}

// Error logs a message at error priority
func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

// Infof logs a message at info priority
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args))
}
