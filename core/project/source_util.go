package project

import (
	"os"
	"strings"

	"kerosenelabs.com/espresso/core/config"
	"kerosenelabs.com/espresso/core/util"
)

// GetProjectPath gets the source path from the given ProjectConfig
func GetSourcePath(cfg *config.ProjectConfig) (*string, error) {
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
