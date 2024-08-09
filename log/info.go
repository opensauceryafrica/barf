package logger

import (
	"log"

	"github.com/opensaucerer/barf/constant"
)

// Info logs a message with the info color
func Info(msg string) {
	log.Println(constant.InfoColor + msg + constant.ResetColor)
}
