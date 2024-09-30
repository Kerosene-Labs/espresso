package project

import (
	"os"
	"strings"

	"hlafaille.xyz/espresso/v0/core/config"
	"hlafaille.xyz/espresso/v0/core/util"
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
