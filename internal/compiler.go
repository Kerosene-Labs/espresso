package internal

import (
	"errors"
	"os/exec"
)

// CompileSourceFile compiles the sourcefile with the given project toolchain
func CompileSourceFile(cfg *ProjectConfig, srcFile *SourceFile) error {
	// run the compiler
	command := cfg.Toolchain.Path + "/bin/javac"
	args := []string{}
	if IsDebugMode() {
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
