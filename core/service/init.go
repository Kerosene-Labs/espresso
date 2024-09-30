// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package service

import (
	"fmt"
	"os"

	"kerosenelabs.com/espresso/core/config"
	"kerosenelabs.com/espresso/core/project"
	"kerosenelabs.com/espresso/core/util"
)

// InitializeProject is a service function for initializing a new project
func InitializeProject(name string, basePkg *string) {
	// ensure JAVA_HOME is set
	javaHome, err := util.GetJavaHome()
	if err != nil {
		fmt.Println("JAVA_HOME is not set, do you have Java installed?")
	}

	// ensure a proejct doesn't already exist
	cfgExists, err := config.DoesExist()
	if err != nil {
		fmt.Println("Error occurred while ensuring a config doesn't already exist")
		panic(err)
	} else if *cfgExists {
		fmt.Println("Config already exists")
		os.Exit(1)
	}

	fmt.Printf("Creating '%s'\n", name)

	// create a base config
	cfg := config.ProjectConfig{
		Name: name,
		Version: config.Version{
			Major:  0,
			Minor:  1,
			Patch:  0,
			Hotfix: nil,
		},
		BasePackage: *basePkg,
		Toolchain: config.Toolchain{
			Path: *javaHome,
		},
		Registries:   []config.Registry{{Name: "espresso-registry", Url: "https://github.com/Kerosene-Labs/espresso-registry/archive/refs/heads/main.zip"}},
		Dependencies: []config.Dependency{},
	}

	// write some example code
	println("Creating base package, writing example code")
	project.WriteExampleCode(&cfg)

	// persist the config
	println("Persisting project configuration")
	config.Persist(&cfg)

	println("Done.")
}
