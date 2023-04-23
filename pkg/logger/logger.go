package logger

import (
	"go.uber.org/zap"
)

// type Log interface {
// 	Info(format string, args ...interface{})
// 	Error(format string, args ...interface{})
// 	Debug(format string, args ...interface{})
// 	Panic(format string, args ...interface{})
// 	Warn(format string, args ...interface{})
// }

// Logger wrapper for logging libraries
type Logger struct {
	logger *zap.SugaredLogger
}

// NewLogger constructs a new Logger class
func NewLogger() *Logger {
	// TODO Add a way to switch between dev vs prod in config
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	defer logger.Sync() // flushes buffer, if any

	return &Logger{sugar}
}

func (l Logger) Info(format string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info(format)
	} else {
		l.logger.Infof(format, args...)
	}
}

func (l Logger) Error(format string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Error(format)
	} else {
		l.logger.Errorf(format, args...)
	}
}

func (l Logger) Debug(format string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Debug(format)
	} else {
		l.logger.Debugf(format, args...)
	}
}

func (l Logger) Panic(format string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Panic(format)
	} else {
		l.logger.Panicf(format, args...)
	}
}

func (l Logger) Warn(format string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Warn(format)
	} else {
		l.logger.Warnf(format, args...)
	}
}
