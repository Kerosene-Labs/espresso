package internal

import (
	"fmt"
	"os"
	"os/exec"
)

// GenerateManifest generates a JVM manifest
func GenerateManifest(cfg *ProjectConfig) string {
	base := "Manifest-Version: 1.0\n"
	base += "Main-Class: " + cfg.BasePackage + ".Main\n"
	base += "Created-By: Espresso"
	return base
}

// Write the Manifest to the build directory
func WriteManifest(cfg *ProjectConfig) error {
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
func PackageClasses(cfg *ProjectConfig) error {
	command := cfg.Toolchain.Path + "/bin/jar"
	args := []string{"cfm"}

	// handle jar output path
	if IsDebugMode() {
		args = append(args, "ESPRESSO_DEBUG/dist/dist.jar")
	} else {
		args = append(args, "dist/dist.jar")
	}

	// write the manifest, include it
	WriteManifest(cfg)
	if IsDebugMode() {
		args = append(args, "ESPRESSO_DEBUG/build/MANIFEST.MF")
	} else {
		args = append(args, "build/MANIFEST.MF")
	}

	// add the class directory
	if IsDebugMode() {
		args = append(args, "-C", "ESPRESSO_DEBUG/build")
	} else {
		args = append(args, "-C", "build")
	}
	args = append(args, ".")

	// run the command
	fmt.Printf("Running: %s %s\n", command, args)
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s", output)
		return err
	}

	return nil
}
