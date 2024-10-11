// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package service

import (
	"fmt"

	"kerosenelabs.com/espresso/core/context/project"
	"kerosenelabs.com/espresso/core/util"
)

// InitializeProject is a service function for initializing a new project
func InitializeProject(name string, basePkg *string) {
	// ensure JAVA_HOME is set
	javaHome, err := util.GetJavaHome()
	if err != nil {
		util.ErrorQuit("JAVA_HOME is not set; do you have java installed?")
	}

	// ensure a proejct doesn't already exist
	// cfgExists, err := project.DoesExist()
	// if err != nil {
	// 	fmt.Println("Error occurred while ensuring a config doesn't already exist")
	// 	panic(err)
	// } else if *cfgExists {
	// 	fmt.Println("Config already exists")
	// 	os.Exit(1)
	// }

	fmt.Printf("Creating '%s'\n", name)

	// create a base config
	config := project.ProjectConfig{
		Name: name,
		Version: project.Version{
			Major:  0,
			Minor:  1,
			Patch:  0,
			Hotfix: nil,
		},
		BasePackage: *basePkg,
		Toolchain: project.Toolchain{
			Path: *javaHome,
		},
		Registries:   []project.Registry{{Name: "espresso-registry", Url: "https://github.com/Kerosene-Labs/espresso-registry/archive/refs/heads/main.zip"}},
		Dependencies: []project.Dependency{},
	}

	// get our config path
	configPath, err := project.GetConfigPath()
	if err != nil {
		util.ErrorQuit("An error occurred while getting the config path: %s", err)
	}

	// create our project context
	projectContext := project.ProjectContext{
		Config:     config,
		ConfigPath: configPath,
	}

	// write some example code
	println("Creating base package, writing example code")
	project.WriteExampleCode(projectContext)

	// persist the config
	println("Persisting project configuration")
	project.Persist(projectContext)

	println("Done.")
}
