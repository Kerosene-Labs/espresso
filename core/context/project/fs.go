package project

import (
	"os"
	"strings"

	"kerosenelabs.com/espresso/core/util"
)

// getConfigPath gets the absolute path to the config file
func GetConfigPath() (string, error) {
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

// getSourcePath gets the path at which there should be source files
func getSourcePath(projectConfig ProjectConfig) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path := wd

	if util.IsDebugMode() {
		path += "/ESPRESSO_DEBUG"
	}

	path += "/src/java/" + strings.ReplaceAll(projectConfig.BasePackage, ".", "/")
	return path, nil
}

// Persist persists the configuration under the EnvironmentContext to the filesystem
func Persist(projectContext ProjectContext) error {
	marshalData, err := MarshalConfig(projectContext.Config)
	if err != nil {
		return err
	}

	// write the config
	cfgExists, err := util.DoesPathExist(projectContext.ConfigPath)
	if err != nil {
		return err
	}

	if cfgExists {
		file, err := os.Open(projectContext.ConfigPath)
		if err != nil {
			return err
		}
		file.WriteString(*marshalData)
		defer file.Close()

	} else {
		file, err := os.Create(projectContext.ConfigPath)
		if err != nil {
			return err
		}
		file.WriteString(*marshalData)
		defer file.Close()
	}

	return nil
}

// WriteExampleCode is an internal function that writes example code to a newly created project
func WriteExampleCode(projectContext ProjectContext) error {
	os.MkdirAll(projectContext.SourcePath, 0755)

	// create the main file
	file, err := os.Create(projectContext.SourcePath + "/Main.java")
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
	code = strings.ReplaceAll(code, "${BASE_PACKAGE}", projectContext.Config.BasePackage)
	file.WriteString(code)
	return nil

}
