package config

import (
	"os"

	"kerosenelabs.com/espresso/core/util"
)

// GetConfigPath gets the absolute path to the config file
func GetConfigPath() (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if util.IsDebugMode() {
		path := wd + "/ESPRESSO_DEBUG" + "/espresso.yml"
		return &path, nil
	} else {
		path := wd + "/espresso.yml"
		return &path, nil
	}
}
