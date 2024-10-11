// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package service

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"kerosenelabs.com/espresso/core/context/project"
	"kerosenelabs.com/espresso/core/dependency"
	"kerosenelabs.com/espresso/core/registry"
	"kerosenelabs.com/espresso/core/source"
	"kerosenelabs.com/espresso/core/toolchain"
	"kerosenelabs.com/espresso/core/util"
)

// BuildProject is a service function for building the current project
func BuildProject() {
	// get our project context
	projectContext, err := project.GetProjectContext()
	if err != nil {
		util.ErrorQuit("An error occurred while getting the project context: %s", err)
	}

	color.Cyan("-- [%s] Beginning build, please ensure you are compliant with all dependency licenses\n", projectContext.Config.Name)

	// discover source files
	files, err := source.DiscoverSourceFiles(projectContext.Config)
	if err != nil {
		util.ErrorQuit(fmt.Sprintf("An error occurred while discovering source files: %s\n", err))
	}

	// run the compiler on each source file
	color.Cyan("-- Compiling")
	var wg sync.WaitGroup
	for _, value := range files {
		wg.Add(1)
		go func(f *source.SourceFile) {
			defer wg.Done()
			err := toolchain.CompileSourceFile(projectContext.Config, value)
			if err != nil {
				util.ErrorQuit("An error occurred while compiling a source file: %s\n", err)
			}
			color.Blue("Compiled: " + f.Path)
		}(&value)
	}
	wg.Wait()

	// package the project
	color.Cyan("-- Packaging distributable")
	err = toolchain.PackageClasses(projectContext.Config)
	if err != nil {
		util.ErrorQuit(fmt.Sprintf("An error occurred while packaging the classes: %s\n", err))
	}
	color.Blue("Finished packaging distributable")

	// iterate over each dependency, resolve it and copy it
	distPath, err := toolchain.GetDistPath(projectContext.Config)
	if err != nil {
		util.ErrorQuit(fmt.Sprintf("Unable to get dist path: %s", err))
	}
	os.MkdirAll(*distPath+"/libs", 0755)
	var depCopyWg sync.WaitGroup
	color.Cyan("-- Copying dependency packages to distributable")
	for _, dep := range projectContext.Config.Dependencies {
		depCopyWg.Add(1)
		go func() {
			defer depCopyWg.Done()
			color.Cyan("-- Beginning copy of '%s:%s' to distributable", dep.Group, dep.Name)
			resolved, err := dependency.ResolveDependency(dep, projectContext.Config.Registries)
			if err != nil {
				util.ErrorQuit(fmt.Sprintf("Unable to resolve dependency: %s", err))
			}

			// calculate the should-be location of this jar locally
			espressoPath, err := util.GetEspressoDirectoryPath()
			if err != nil {
				util.ErrorQuit(fmt.Sprintf("Unable to get the espresso home: %s", espressoPath))
			}
			signature := registry.CalculatePackageSignature(resolved.Package, resolved.PackageVersion)
			cachedPackageHome := espressoPath + "/cachedPackages" + signature + ".jar"

			// copy the file
			util.CopyFile(cachedPackageHome, *distPath+"/libs")

			color.Blue("Copied '%s:%s' to distributable", dep.Group, dep.Name)
		}()
	}
	depCopyWg.Wait()
	color.Green("Done!")
}
