// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package toolchain

import (
	"errors"
	"os"
	"os/exec"

	"strings"

	"kerosenelabs.com/espresso/core/context/project"
	"kerosenelabs.com/espresso/core/dependency"
	"kerosenelabs.com/espresso/core/util"
)

func capLinesAt72Bytes(input string) []string {
	var lines []string
	var currentLine string
	var currentLineBytes int

	words := strings.Fields(input) // Split the input into words based on spaces

	for _, word := range words {
		// Prepare the word to be added, including a preceding space if necessary
		wordToAdd := word
		if currentLineBytes > 0 {
			wordToAdd = " " + word
		}
		wordBytes := len(wordToAdd)

		if currentLineBytes+wordBytes > 72 {
			// Line length exceeded; append the current line to lines
			lines = append(lines, currentLine)
			// Start a new line with a leading space as per manifest continuation rules
			currentLine = " " + word
			currentLineBytes = len(currentLine)
		} else {
			// Append the word to the current line
			currentLine += wordToAdd
			currentLineBytes += wordBytes
		}
	}

	// Append the last line if it's not empty
	if currentLineBytes > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

// GenerateManifest generates a JVM manifest
func GenerateManifest(cfg project.ProjectConfig) (string, error) {
	base := "Manifest-Version: 1.0\n"
	base += "Main-Class: " + cfg.BasePackage + ".Main\n"
	base += "Created-By: Espresso\n"

	// iterate over dependencies, resolve them and add them to the manifest base
	if len(cfg.Dependencies) > 0 {
		classPath := "Class-Path: "
		for _, dep := range cfg.Dependencies {
			resolvedDependency, err := dependency.ResolveDependency(dep, cfg.Registries)
			if err != nil {
				return "", err
			}
			classPath += "libs/" + resolvedDependency.Package.Name + ".jar "
		}

		splitLines := capLinesAt72Bytes(classPath)
		for _, line := range splitLines {
			base += line + "\n"
		}
	}

	println(strings.TrimSuffix(base, "\n"))
	return strings.TrimSuffix(base, "\n"), nil
}

// Write the Manifest to the build directory
func WriteManifest(cfg project.ProjectConfig) error {
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
	content, err := GenerateManifest(cfg)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

// PackageClasses creates a .jar of the given classes
func PackageClasses(cfg project.ProjectConfig) error {
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
