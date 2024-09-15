package toolchain

import (
	"errors"
	"os"
)

// GetJavaHome gets the value of the JAVA_HOME
func GetJavaHome() (*string, error) {
	path, present := os.LookupEnv("JAVA_HOME")
	if !present {
		return nil, errors.New("java_home is not set")
	}
	return &path, nil
}
