package dependency

import (
	"errors"
	"fmt"

	"kerosenelabs.com/espresso/core/config"
	"kerosenelabs.com/espresso/core/project"
	"kerosenelabs.com/espresso/core/registry"
)

// ResolvedDependency represents a match between a project dependency and a registry package.
type ResolvedDependency struct {
	Dependency     *config.Dependency
	Package        *registry.Package
	PackageVersion *registry.PackageVersionDeclaration
}

// ResolveDependency resolves the given dependency. This function will iterate over the given registries,
// their packages, and each version of that package. The first match is the once that will be returned.
// Registries follow hierarchical order, so the top-most one is the one that is searched first.
func ResolveDependency(dep *config.Dependency, regs *[]config.Registry) (*ResolvedDependency, error) {
	if dep == nil || regs == nil {
		return nil, errors.New("required argument was a nil pointer")
	}

	// get our version string
	depVersionStr := project.GetVersionAsString(&dep.Version)

	// iterate over each registry
	for _, reg := range *regs {
		// get this registry's packages on the filesystem cache
		pkgs, err := registry.GetRegistryPackages(reg)
		if err != nil {
			return nil, err
		}
		// iterate over each package
		for _, pkg := range pkgs {
			// if we have a match via group and name, match a version
			if pkg.Group == dep.Group && pkg.Name == dep.Name {
				for _, version := range pkg.Versions {
					if version.Number == depVersionStr {
						return &ResolvedDependency{
							Dependency:     dep,
							Package:        &pkg,
							PackageVersion: &version,
						}, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("'%s:%s:%s' dependency was unable to be resolved within any given registry", dep.Group, dep.Name, depVersionStr)
}
