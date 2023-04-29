package logger

import (
	"log"

	"github.com/opensaucerer/barf/config"
)

// Error logs an error message
func Error(msg string) {
	log.Println(config.ErrorColor + msg + config.ResetColor)
}
