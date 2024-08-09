package logger

import (
	"log"

	"github.com/opensaucerer/barf/constant"
)

// Warn prints a warning message
func Warn(msg string) {
	log.Println(constant.WarnColor + msg + constant.ResetColor)
}
