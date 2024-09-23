package dependency

import (
	"io/fs"
	"path/filepath"

	"hlafaille.xyz/espresso/v0/core/project"
)

// getPackageGroupPathsInRegCache walks over a particular registry's lookup directory, (ex: lookup/espresso-registry-main)
// and looks for group directories (ex: org.projectlombok)
func getPackageGroupPathsInRegCache(reg project.Registry) ([]string, error) {
	// get the cache path
	cachePath, err := GetCachePath(&reg)
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

// getPackageDeclarationPathsInRegistryCache walks over a particular package group directory
// in the registry's lookup (ex: lookup/espresso-registry-main/org.projectlombok) and looks for package
// declaration .yml files (ex: lombok.yml), returning the paths to each one of those said .yml files.
func getPackageDeclarationPathsInRegistryCache(reg *project.Registry, pkgGroupPath *string) ([]string, error) {
	if reg == nil {
		panic("programming error: reg is nil")
	}

	// get the cache path
	cachePath, err := GetCachePath(reg)
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
