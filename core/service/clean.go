// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package service

import (
	"os"

	"github.com/fatih/color"
	"kerosenelabs.com/espresso/core/context/project"
	"kerosenelabs.com/espresso/core/toolchain"
	"kerosenelabs.com/espresso/core/util"
)

// CleanWorkspace is a service function to clean the current workspace.
func CleanWorkspace() {
	// get our project context
	projectContext, err := project.GetProjectContext()
	if err != nil {
		util.ErrorQuit("An error occurred while getting the project context: %s", err)
	}

	// get the build dir
	buildPath, err := toolchain.GetBuildPath(projectContext.Config)
	if err != nil {
		util.ErrorQuit("An error occurred while getting the build path: %s", err)
	}

	// get the dist dir
	distPath, err := toolchain.GetDistPath(projectContext.Config)
	if err != nil {
		util.ErrorQuit("An error occurred while getting the build path: %s", err)
	}

	// remove the build dir
	err = os.RemoveAll(*buildPath)
	if err != nil {
		util.ErrorQuit("An error occurred while deleting the build path: %s", err)
	}

	// remove the dist dir
	err = os.RemoveAll(*distPath)
	if err != nil {
		util.ErrorQuit("An error occurred while deleting the dist path: %s", err)
	}

	color.Green("Cleaned")
}
