package registry

import (
	"io/fs"
	"os"
	"path/filepath"

	"hlafaille.xyz/espresso/v0/core/project"
)

// GetCachePath gets the full cache path from the registry (ex: /home/vscode/.espresso/registries/espresso-registry)
func GetRegistryCachePath(reg *project.Registry) (string, error) {
	// get our home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.espresso/registries/" + reg.Name, nil
}

// GetRegistryCacheDependenciesPath
func GetRegistryCacheDependenciesPath(reg *project.Registry) (string, error) {
	cachePath, err := GetRegistryCachePath(reg)
	if err != nil {
		return "", err
	}
	return cachePath + "/espresso-registry-main/dependencies", nil
}

// walkRegistryLookup walks over a particular registry's lookup directory, (ex: lookup/espresso-registry-main)
// and looks for group directories (ex: org.projectlombok)
func walkRegistryLookup(reg project.Registry) ([]string, error) {
	// get the cache path
	cachePath, err := GetRegistryCachePath(&reg)
	if err != nil {
		return []string{}, err
	}

	// walk the directory for all groupId's
	var dirs []string = []string{}
	err = filepath.Walk(cachePath+"/lookup/espresso-registry-main/dependencies", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "dependencies" {
			return nil
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	return dirs, nil
}

// walkPackageGroup walks a package group for package declaration files (.yml). Returns a list of their paths.
func walkPackageGroup(packageGroupPath string) ([]string, error) {
	var paths []string = []string{}
	err := filepath.Walk(packageGroupPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return paths, nil
}
