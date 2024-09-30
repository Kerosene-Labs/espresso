package project

import (
	"fmt"
	"os"
	"strings"

	"kerosenelabs.com/espresso/core/config"
)

// WriteExampleCode is an internal function that writes example code to a newly created project
func WriteExampleCode(cfg *config.ProjectConfig) error {
	path, _ := GetSourcePath(cfg)
	os.MkdirAll(*path, 0755)

	// create the main file
	file, err := os.Create(*path + "/Main.java")
	if err != nil {
		return err
	}
	defer file.Close()

	// write some code
	code := `package ${BASE_PACKAGE};

import java.lang.System;

public class Main {
    public static void main(String[] args) {
        System.out.println("Hello from Espresso!");
    }
}`
	code = strings.ReplaceAll(code, "${BASE_PACKAGE}", cfg.BasePackage)
	file.WriteString(code)
	return nil

}

// GetConfigPath gets the absolute path to the config file
// func GetConfigPath() (*string, error) {
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if util.IsDebugMode() {
// 		path := wd + "/ESPRESSO_DEBUG" + "/espresso.yml"
// 		return &path, nil
// 	} else {
// 		path := wd + "/espresso.yml"
// 		return &path, nil
// 	}
// }

// GetVersionAsString gets a project version as a string.
func GetVersionAsString(version *config.Version) string {
	if version == nil {
		panic("version is nil")
	}
	versionString := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)
	if version.Hotfix != nil {
		versionString += *version.Hotfix
	}
	return versionString
}
