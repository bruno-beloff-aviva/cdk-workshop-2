package loglevel

import (
	"github.com/labstack/gommon/log"
)

func Level(description string) log.Lvl {
	switch description {
	case "DEBUG":
		return log.DEBUG
	case "INFO":
		return log.INFO
	case "WARN":
		return log.WARN
	case "ERROR":
		return log.ERROR
	case "OFF":
		return log.OFF

	default:
		return log.INFO
	}
}
