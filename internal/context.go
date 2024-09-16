package internal

import "os"

// IsDebugMode returns if we should treat the current runtime context and all actions as debug mode.
// Debug mode is effectively wrapping all filesystem actions in the "espresso_debug" directory.
func IsDebugMode() bool {
	val, present := os.LookupEnv("ESPRESSO_DEBUG")
	if !present {
		return false
	}
	return val == "1"
}
