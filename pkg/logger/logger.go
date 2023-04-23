package logger

import (
	"log"

	"go.uber.org/zap"
)

// type Log interface {
// 	Info(format string, args ...interface{})
// 	Error(format string, args ...interface{})
// 	Debug(format string, args ...interface{})
// 	Panic(format string, args ...interface{})
// 	Warn(format string, args ...interface{})
// }

var globalLogger *logger

// TODO Add a way to switch between dev vs prod in config
// InitLogging initializes the global logger
func InitLogging() error {
	// Setup the zapper framework
	zapper, err := zap.NewDevelopment()
	if err != nil {
		log.Panicf("Failed to create zapper logger", err)
	}
	sugar := zapper.Sugar()
	defer zapper.Sync() // flushes buffer, if any

	globalLogger = &logger{sugar}
	Info("logger created")
	return nil
}

// logger is a wrapper around the zap logger
type logger struct {
	log *zap.SugaredLogger
}

func Info(format string, args ...interface{}) {
	if len(args) == 0 {
		globalLogger.log.Info(format)
	} else {
		globalLogger.log.Infof(format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if len(args) == 0 {
		globalLogger.log.Error(format)
	} else {
		globalLogger.log.Errorf(format, args...)
	}
}

func Debug(format string, args ...interface{}) {
	if len(args) == 0 {
		globalLogger.log.Debug(format)
	} else {
		globalLogger.log.Debugf(format, args...)
	}
}

func Panic(format string, args ...interface{}) {
	if len(args) == 0 {
		globalLogger.log.Panic(format)
	} else {
		globalLogger.log.Panicf(format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	if len(args) == 0 {
		globalLogger.log.Warn(format)
	} else {
		globalLogger.log.Warnf(format, args...)
	}
}

func Fatal(format string, args ...interface{}) {
	if len(args) == 0 {
		globalLogger.log.Fatal(format)
	} else {
		globalLogger.log.Fatalf(format, args...)
	}
}
