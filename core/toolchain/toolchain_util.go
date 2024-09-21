package toolchain

import (
	"os"

	"hlafaille.xyz/espresso/v0/core/project"
	"hlafaille.xyz/espresso/v0/core/util"
)

// GetBuildPath gets the absolute path to the build directory
func GetBuildPath(cfg *project.ProjectConfig) (*string, error) {
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
func GetDistPath(cfg *project.ProjectConfig) (*string, error) {
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
