// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package source

import (
	"os"
	"strings"

	"kerosenelabs.com/espresso/core/context/project"
	"kerosenelabs.com/espresso/core/util"
)

// GetProjectPath gets the source path from the given ProjectConfig
func GetSourcePath(cfg project.ProjectConfig) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path := wd

	if util.IsDebugMode() {
		path += "/ESPRESSO_DEBUG"
	}

	path += "/src/java/" + strings.ReplaceAll(cfg.BasePackage, ".", "/")
	return &path, nil
}
