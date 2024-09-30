// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package config

import "gopkg.in/yaml.v3"

// Dependency represents a particular dependency
type Dependency struct {
	Group   string  `yaml:"group"`
	Name    string  `yaml:"name"`
	Version Version `yaml:"version"`
}

// Registry represents a particular repository
type Registry struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
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
	Dependencies []Dependency `yaml:"dependencies"`
	Registries   []Registry   `yaml:"registries"`
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
