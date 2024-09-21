package project

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"hlafaille.xyz/espresso/v0/core/util"
)

// Dependency represents a particular dependency
type Dependency struct {
	GroupID    string `yaml:"groupId"`
	ArtifactID string `yaml:"artifactId"`
	Version    string `yaml:"version"`
}

// Registry represents a particular repository
type Registry struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

// Dependencies represents dependency management configuration
type Dependencies struct {
	Registries   []Registry   `yaml:"registries"`
	Dependencies []Dependency `yaml:"uses"`
}

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

// ProjectConfig represents an Espresso project
type ProjectConfig struct {
	Name         string       `yaml:"name"`
	Version      Version      `yaml:"version"`
	BasePackage  string       `yaml:"basePackage"`
	Toolchain    Toolchain    `yaml:"toolchain"`
	Dependencies Dependencies `yaml:"dependencies"`
}

// UnmarshalConfig marshals the given ProjectConfig to yml
func UnmarshalConfig(cfgYml string) (*ProjectConfig, error) {
	var cfg ProjectConfig
	err := yaml.Unmarshal([]byte(cfgYml), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
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

// GetConfig reads the config from the filesystem and returns a pointer to a ProjectConfig, or an error
func GetConfig() (*ProjectConfig, error) {
	// get the config path for this context
	path, err := util.GetConfigPath()
	if err != nil {
		return nil, err
	}

	// read the file
	file, err := os.ReadFile(*path)
	if err != nil {
		return nil, err
	}

	// unmarshal into a cfg struct
	cfg, err := UnmarshalConfig(string(file))
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// ConfigExists gets if a config exists at the current directory
func ConfigExists() (*bool, error) {
	configPath, err := util.GetConfigPath()
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

	cfgPath, err := util.GetConfigPath()
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
