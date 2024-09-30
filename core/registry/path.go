package registry

import (
	"io/fs"
	"path/filepath"

	"hlafaille.xyz/espresso/v0/core/project"
	"hlafaille.xyz/espresso/v0/core/util"
)

// GetCachePath gets the full cache path from the registry (ex: /home/vscode/.espresso/registries/espresso-registry)
func GetRegistryCachePath(reg *project.Registry) (string, error) {
	// get our home dir
	espressoPath, err := util.GetEspressoDirectoryPath()
	if err != nil {
		return "", err
	}
	return espressoPath + "/registries/" + reg.Name, nil
}

// GetRegistryCachePackagesLookupPath gets the full cache path of the registry package lookup
// (ex: /home/vscode/.espresso/registries/espresso-registry/lookup/espresso-registry-main/packages)
func GetRegistryCachePackagesLookupPath(reg *project.Registry) (string, error) {
	regCachePath, err := GetRegistryCachePath(reg)
	if err != nil {
		return "", err
	}
	return regCachePath + "/lookup/espresso-registry-main/packages", nil
}

// walkRegistryLookup walks over a particular registry's lookup directory, (ex: lookup/espresso-registry-main)
// and looks for group directories (ex: org.projectlombok)
func walkRegistryLookup(reg project.Registry) ([]string, error) {
	// get the cache path
	cachePath, err := GetRegistryCachePackagesLookupPath(&reg)
	if err != nil {
		return []string{}, err
	}

	// walk the directory for all groupId's
	var dirs []string = []string{}
	err = filepath.Walk(cachePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "packages" {
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
