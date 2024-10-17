// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package dependency

import (
	"fmt"

	"kerosenelabs.com/espresso/core/context/project"
	"kerosenelabs.com/espresso/core/registry"
	"kerosenelabs.com/espresso/core/util"
)

// ResolvedDependency represents a match between a project dependency and a registry package.
type ResolvedDependency struct {
	Dependency       project.Dependency
	Package          registry.Package
	PackageVersion   registry.PackageVersionDeclaration
	Registry         project.Registry
	PackageSignature string
}

// GetCachePath returns the absolute filesystem path to where the cached package should be located. Please note,
// this function will always return the path of the .jar even if it does not exist.
//
// You should use PackageCachePath.DoesExist() to ensure this path exists if you're depennding on the package.
func (resolvedDependency ResolvedDependency) GetCachePath() (util.Path, error) {
	packageSignature := registry.CalculatePackageSignature(resolvedDependency.Registry, resolvedDependency.Package, resolvedDependency.PackageVersion)
	espressoPath, err := util.GetEspressoDirectoryPath()
	if err != nil {
		return util.Path{}, err
	}
	pkgPath := espressoPath + "/cachedPackages/" + packageSignature + ".jar"
	return util.Path{Absolute: pkgPath}, nil
}

// ResolveDependency resolves the given dependency. This function will iterate over the given registries,
// their packages, and each version of that package. The first match is the once that will be returned.
// Registries follow hierarchical order, so the top-most one is the one that is searched first.
func ResolveDependency(dependency project.Dependency, registries []project.Registry) (ResolvedDependency, error) {
	// get our version string
	depVersionStr := project.GetVersionAsString(dependency.Version)

	// iterate over each registry
	for _, reg := range registries {
		// get this registry's packages on the filesystem cache
		pkgs, err := registry.GetRegistryPackages(reg)
		if err != nil {
			return ResolvedDependency{}, err
		}
		// iterate over each package
		for _, pkg := range pkgs {
			// if we have a match via group and name, match a version
			if pkg.Group == dependency.Group && pkg.Name == dependency.Name {
				for _, version := range pkg.Versions {
					if version.Number == depVersionStr {
						return ResolvedDependency{
							Dependency:       dependency,
							Package:          pkg,
							PackageVersion:   version,
							Registry:         reg,
							PackageSignature: registry.CalculatePackageSignature(reg, pkg, version),
						}, nil
					}
				}
			}
		}
	}
	return ResolvedDependency{}, fmt.Errorf("'%s:%s:%s' dependency was unable to be resolved within any given registry", dependency.Group, dependency.Name, depVersionStr)
}
