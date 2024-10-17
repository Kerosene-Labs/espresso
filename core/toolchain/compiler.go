// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package toolchain

import (
	"errors"
	"os/exec"

	"kerosenelabs.com/espresso/core/context/project"
	"kerosenelabs.com/espresso/core/dependency"
	"kerosenelabs.com/espresso/core/source"
	"kerosenelabs.com/espresso/core/util"
)

// CompileSourceFile compiles the sourcefile with the given project toolchain
func CompileSourceFile(cfg project.ProjectConfig, srcFile source.SourceFile) error {
	// initialize our classpath value
	cpVal := ""
	if util.IsDebugMode() {
		cpVal = "ESPRESSO_DEBUG/src/java"
	} else {
		cpVal = "src/java"
	}

	// iterate over dependencies, resolve each one and add it to the classpath argument value
	for _, dep := range cfg.Dependencies {
		// resolve our dependency
		resolvedDependency, err := dependency.ResolveDependency(dep, cfg.Registries)
		if err != nil {
			return err
		}

		// get our cache path for the jar
		depCachePath, err := resolvedDependency.GetCachePath()
		if err != nil {
			return err
		}

		// append it to the classpath value
		cpVal += ":" + depCachePath.Absolute
	}

	// run the compiler
	command := cfg.Toolchain.Path + "/bin/javac"
	args := []string{}
	if util.IsDebugMode() {
		args = append(args, "-cp", cpVal, "-d", "ESPRESSO_DEBUG/build")
	} else {
		args = append(args, "-cp", cpVal, "-d", "build")
	}
	args = append(args, srcFile.Path)
	cmd := exec.Command(command, args...)

	// handle output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}
