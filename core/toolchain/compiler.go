package toolchain

import (
	"errors"
	"os/exec"

	"hlafaille.xyz/espresso/v0/core/project"
	"hlafaille.xyz/espresso/v0/core/util"
)

// CompileSourceFile compiles the sourcefile with the given project toolchain
func CompileSourceFile(cfg *project.ProjectConfig, srcFile *project.SourceFile) error {
	// run the compiler
	command := cfg.Toolchain.Path + "/bin/javac"
	args := []string{}
	if util.IsDebugMode() {
		args = append(args, "-cp", "ESPRESSO_DEBUG/src/java", "-d", "ESPRESSO_DEBUG/build")
	} else {
		args = append(args, "-cp", "src/java", "build")
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
