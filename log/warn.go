package logger

import (
	"log"

	"github.com/opensaucerer/barf/config"
)

// Warn prints a warning message
func Warn(msg string) {
	log.Println(config.WarnColor + msg + config.ResetColor)
}
