package toolchain

import (
	"errors"
	"os"
	"os/exec"

	"hlafaille.xyz/espresso/v0/core/config"
	"hlafaille.xyz/espresso/v0/core/util"
)

// GenerateManifest generates a JVM manifest
func GenerateManifest(cfg *config.ProjectConfig) string {
	base := "Manifest-Version: 1.0\n"
	base += "Main-Class: " + cfg.BasePackage + ".Main\n"
	base += "Created-By: Espresso"
	return base
}

// Write the Manifest to the build directory
func WriteManifest(cfg *config.ProjectConfig) error {
	// get the path where it should live
	buildPath, err := GetBuildPath(cfg)
	path := *buildPath + "/MANIFEST.MF"
	if err != nil {
		return err
	}

	// open the file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// write the file
	content := GenerateManifest(cfg)
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

// PackageClasses creates a .jar of the given classes
func PackageClasses(cfg *config.ProjectConfig) error {
	command := cfg.Toolchain.Path + "/bin/jar"
	args := []string{"cfm"}

	// handle jar output path
	if util.IsDebugMode() {
		args = append(args, "ESPRESSO_DEBUG/dist/dist.jar")
	} else {
		args = append(args, "dist/dist.jar")
	}

	// write the manifest, include it
	WriteManifest(cfg)
	if util.IsDebugMode() {
		args = append(args, "ESPRESSO_DEBUG/build/MANIFEST.MF")
	} else {
		args = append(args, "build/MANIFEST.MF")
	}

	// add the class directory
	if util.IsDebugMode() {
		args = append(args, "-C", "ESPRESSO_DEBUG/build")
	} else {
		args = append(args, "-C", "build")
	}
	args = append(args, ".")

	// run the command
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	return nil
}
