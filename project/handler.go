package project

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Toolchain represents the toolchain on the system
type Toolchain struct {
	Path string `yaml:"path"`
}

// ProjectVersion represents a semantic version number
type Version struct {
	Major  int64   `yaml:"major"`
	Minor  int64   `yaml:"minor"`
	Patch  int64   `yaml:"patch"`
	Hotfix *string `yaml:"hotfix"`
}

// Project represents an Espresso project
type ProjectConfig struct {
	Name        string    `yaml:"name"`
	Version     Version   `yaml:"version"`
	BasePackage string    `yaml:"base_package"`
	Toolchain   Toolchain `yaml:"toolchain"`
}

// MarshalConfig marshals the given ProjectConfig to yml
func MarshalConfig(cfg *ProjectConfig) (*string, error) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	resp := string(data)
	return &resp, nil
}

// ConfigExists gets if a config exists at the current directory
func ConfigExists() (*bool, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(*configPath)
	resp := err == nil || !os.IsNotExist(err)
	return &resp, nil
}

// WriteExampleCode is an internal function that writes example code to a newly created project
func WriteExampleCode(cfg *ProjectConfig) error {
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

// PersistConfig persists the given ProjectConfig to the filesystem
func PersistConfig(cfg *ProjectConfig) error {
	marshalData, err := MarshalConfig(cfg)
	if err != nil {
		return err
	}

	cfgPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// write the config
	cfgExists, err := ConfigExists()
	if err != nil {
		return err
	}

	if *cfgExists {
		file, err := os.Open(*cfgPath)
		if err != nil {
			return err
		}
		file.WriteString(*marshalData)
		defer file.Close()

	} else {
		file, err := os.Create(*cfgPath)
		if err != nil {
			return err
		}
		file.WriteString(*marshalData)
		defer file.Close()
	}

	return nil
}
