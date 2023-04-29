package logger

import (
	"log"

	"github.com/opensaucerer/barf/config"
)

// Info logs a message with the info color
func Info(msg string) {
	log.Println(config.InfoColor + msg + config.ResetColor)
}
