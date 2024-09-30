package cli

import (
	"os"

	"github.com/fatih/color"
)

func ErrorQuit(msg string) {
	color.Red(msg)
	os.Exit(1)
}
