package logger

import (
	"log"

	"github.com/opensaucerer/barf/constant"
)

// Debug logs a debug message
func Debug(msg string) {
	log.Println(constant.DebugColor + msg + constant.ResetColor)
}
