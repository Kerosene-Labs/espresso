package util

import (
	"os"

	"github.com/fatih/color"
)

func ErrorQuit(msg string, args ...any) {
	color.Red(msg+"\n", args...)
	os.Exit(1)
}
