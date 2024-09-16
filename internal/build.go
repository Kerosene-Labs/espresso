package internal

import (
	"os/exec"
)

// CompileSourceFile compiles the sourcefile with the given project toolchain
func CompileSourceFile(cfg *ProjectConfig, srcFile *SourceFile) error {
	// run the compiler
	command := cfg.Toolchain.Path + "/bin/javac"
	args := []string{}
	if IsDebugMode() {
		args = append(args, "-d", "ESPRESSO_DEBUG/build")
	} else {
		args = append(args, "build")
	}
	args = append(args, srcFile.Path)
	cmd := exec.Command(command, args...)

	// handle output
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
