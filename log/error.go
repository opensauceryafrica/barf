package logger

import (
	"log"

	"github.com/opensaucerer/barf/constant"
)

// Error logs an error message
func Error(msg string) {
	log.Println(constant.ErrorColor + msg + constant.ResetColor)
}
