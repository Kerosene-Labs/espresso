package context

import (
	"os"

	"kerosenelabs.com/espresso/core/util"
)

// getConfigPath gets the absolute path to the config file
func getConfigPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if util.IsDebugMode() {
		path := wd + "/ESPRESSO_DEBUG" + "/espresso.yml"
		return path, nil
	} else {
		path := wd + "/espresso.yml"
		return path, nil
	}
}
