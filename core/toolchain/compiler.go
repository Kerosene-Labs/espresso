// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package toolchain

import (
	"errors"
	"os/exec"

	"kerosenelabs.com/espresso/core/config"
	"kerosenelabs.com/espresso/core/project"
	"kerosenelabs.com/espresso/core/util"
)

// CompileSourceFile compiles the sourcefile with the given project toolchain
func CompileSourceFile(cfg *config.ProjectConfig, srcFile *project.SourceFile) error {
	// build our classpath value
	cpVal := ""
	if util.IsDebugMode() {
		cpVal = "ESPRESSO_DEBUG/src/java"
	} else {
		cpVal = "src/java"
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
