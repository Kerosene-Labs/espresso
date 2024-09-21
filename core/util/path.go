package util

import (
	"os"
)

// GetConfigPath gets the absolute path to the config file
func GetConfigPath() (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if IsDebugMode() {
		path := wd + "/ESPRESSO_DEBUG" + "/espresso.yml"
		return &path, nil
	} else {
		path := wd + "/espresso.yml"
		return &path, nil
	}
}
