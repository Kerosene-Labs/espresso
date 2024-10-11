// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package toolchain

import (
	"os"

	"kerosenelabs.com/espresso/core/context/project"
	"kerosenelabs.com/espresso/core/util"
)

// GetBuildPath gets the absolute path to the build directory
func GetBuildPath(cfg project.ProjectConfig) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path := wd

	if util.IsDebugMode() {
		path += "/ESPRESSO_DEBUG"
	}

	path += "/build"
	return &path, nil
}

// GetDistPath gets the absolute path to the dist directory
func GetDistPath(cfg project.ProjectConfig) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path := wd

	if util.IsDebugMode() {
		path += "/ESPRESSO_DEBUG"
	}

	path += "/dist"
	return &path, nil
}
