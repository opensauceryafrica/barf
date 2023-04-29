package logger

import (
	"log"

	"github.com/opensaucerer/barf/config"
)

// Debug logs a debug message
func Debug(msg string) {
	log.Println(config.DebugColor + msg + config.ResetColor)
}
