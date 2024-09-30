package util

import (
	"errors"
	"os"
)

// IsDebugMode returns if we should treat the current runtime context and all actions as debug mode.
// Debug mode is effectively wrapping all filesystem actions in the "espresso_debug" directory.
func IsDebugMode() bool {
	val, present := os.LookupEnv("ESPRESSO_DEBUG")
	if !present {
		return false
	}
	return val == "1"
}

// GetJavaHome gets the value of the JAVA_HOME
func GetJavaHome() (*string, error) {
	path, present := os.LookupEnv("JAVA_HOME")
	if !present {
		return nil, errors.New("java_home is not set")
	}
	return &path, nil
}

func DoesPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	resp := err == nil || !os.IsNotExist(err)
	return resp, nil
}
