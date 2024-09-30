// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

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
