package project

import (
	"os"
	"strings"
)

// IsDebugMode returns if we should treat the current runtime context and all actions as debug mode.
// Debug mode is effectively wrapping all filesystem actions in the "espresso_debug" directory.
func IsDebugMode() bool {
	val, present := os.LookupEnv("ESPRESSO_DEBUG")
	if !present {
		return false
	}
	return val == "1"
}

// GetProjectPath gets the source path from the given ProjectConfig
func GetSourcePath(cfg *ProjectConfig) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path := wd

	if IsDebugMode() {
		path += "/ESPRESSO_DEBUG"
	}

	path += "/src/java/" + strings.ReplaceAll(cfg.BasePackage, ".", "/")
	return &path, nil
}

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

// GetBuildPath gets the absolute path to the build directory
func GetBuildPath(cfg *ProjectConfig) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path := wd

	if IsDebugMode() {
		path += "/ESPRESSO_DEBUG"
	}

	path += "/build"
	return &path, nil
}
