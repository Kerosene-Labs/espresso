// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package util

import (
	"os"

	"github.com/fatih/color"
)

func ErrorQuit(msg string, args ...any) {
	color.Red(msg+"\n", args...)
	os.Exit(1)
}
