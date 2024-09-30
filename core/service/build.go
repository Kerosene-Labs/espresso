package service

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"kerosenelabs.com/espresso/core/context"
	"kerosenelabs.com/espresso/core/dependency"
	"kerosenelabs.com/espresso/core/project"
	"kerosenelabs.com/espresso/core/registry"
	"kerosenelabs.com/espresso/core/toolchain"
	"kerosenelabs.com/espresso/core/util"
)

// BuildProject is a service function for building the current project
func BuildProject() {
	// get the environment context
	ctx, err := context.GetEnvironmentContext()
	if err != nil {
		util.ErrorQuit("An error occurred while getting the environment context: %s", err)
	}

	// get our project context
	prjCtx, err := ctx.GetProjectContext()
	if err != nil {
		util.ErrorQuit("An error occurred while getting the project context: %s", err)
	}

	color.Cyan("-- Building '%s', please ensure you are compliant with all dependency licenses\n", prjCtx.Cfg.Name)

	// discover source files
	files, err := project.DiscoverSourceFiles(prjCtx.Cfg)
	if err != nil {
		util.ErrorQuit(fmt.Sprintf("An error occurred while discovering source files: %s\n", err))
	}

	// run the compiler on each source file
	var wg sync.WaitGroup
	for _, value := range files {
		wg.Add(1)
		go func(f *project.SourceFile) {
			defer wg.Done()
			color.Cyan("-- Compiling: " + f.Path)
			err := toolchain.CompileSourceFile(prjCtx.Cfg, &value)
			if err != nil {
				util.ErrorQuit("An error occurred while compiling a source file: %s\n", err)
			}
		}(&value)
	}
	wg.Wait()

	// package the project
	color.Cyan("-- Packaging distributable")
	err = toolchain.PackageClasses(prjCtx.Cfg)
	if err != nil {
		util.ErrorQuit(fmt.Sprintf("An error occurred while packaging the classes: %s\n", err))
	}
	color.Blue("Finished packaging distributable")

	// iterate over each dependency, resolve it and copy it
	distPath, err := toolchain.GetDistPath(prjCtx.Cfg)
	if err != nil {
		util.ErrorQuit(fmt.Sprintf("Unable to get dist path: %s", err))
	}
	os.MkdirAll(*distPath+"/libs", 0755)
	var depCopyWg sync.WaitGroup
	color.Cyan("-- Copying packages to distributable")
	for _, dep := range prjCtx.Cfg.Dependencies {
		depCopyWg.Add(1)
		go func() {
			defer depCopyWg.Done()
			color.Cyan("-- Beginning copy of '%s:%s' to distributable", dep.Group, dep.Name)
			resolved, err := dependency.ResolveDependency(&dep, &prjCtx.Cfg.Registries)
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
