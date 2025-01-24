package log

import (
	"fmt"

	gommon "github.com/labstack/gommon/log"
)

var Logger *gommon.Logger

func Level(description string) gommon.Lvl {
	switch description {
	case "DEBUG":
		return gommon.DEBUG
	case "INFO":
		return gommon.INFO
	case "WARN":
		return gommon.WARN
	case "ERROR":
		return gommon.ERROR
	case "OFF":
		return gommon.OFF

	default:
		fmt.Println("Unknown log level: ", description)
		return gommon.INFO
	}
}

func Init(name string) { // , level string
	Logger = gommon.New(name)

	// intLevel := Level(level)
	// Logger.SetLevel(intLevel)
}

// func Debug(args ...interface{}) {
// 	logger.Debug(args...)
// }

// func Debugf(format string, args ...interface{}) {
// 	logger.Debugf(format, args...)
// }

// func Info(args ...interface{}) {
// 	logger.Info(args...)
// }

// func Infof(format string, args ...interface{}) {
// 	logger.Infof(format, args...)
// }

// func Warn(args ...interface{}) {
// 	logger.Warn(args...)
// }

// func Warnf(format string, args ...interface{}) {
// 	logger.Warnf(format, args...)
// }

// func Error(args ...interface{}) {
// 	logger.Error(args...)
// }

// func Errorf(format string, args ...interface{}) {
// 	logger.Errorf(format, args...)
// }
