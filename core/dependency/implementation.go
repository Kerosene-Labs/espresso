package dependency

import (
	"errors"
	"fmt"
	"os"

	"hlafaille.xyz/espresso/v0/core/project"
	"hlafaille.xyz/espresso/v0/core/registry"
	"hlafaille.xyz/espresso/v0/core/util"
)

// ResolvedDependency represents a match between a project dependency and a registry package.
type ResolvedDependency struct {
	Dependency     *project.Dependency
	Package        *registry.Package
	PackageVersion *registry.PackageVersionDeclaration
}

// ResolveDependency resolves the given dependency. This function will iterate over the given registries,
// their packages, and each version of that package. The first match is the once that will be returned.
// Registries follow hierarchical order, so the top-most one is the one that is searched first.
func ResolveDependency(dep *project.Dependency, regs *[]project.Registry) (*ResolvedDependency, error) {
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

// CacheResolvedDependency fetches the resolved dependency from the internet
func CacheResolvedDependency(rdep *ResolvedDependency) error {
	if rdep == nil {
		return errors.New("rdep was nil")
	}

	// get our package signature
	packageSignature := registry.CalculatePackageSignature(rdep.Package, rdep.PackageVersion)

	// get where we should store this package
	espressoPath, err := util.GetEspressoDirectoryPath()
	if err != nil {
		return err
	}
	pkgPath := espressoPath + "/cachedPackages/" + packageSignature + ".jar"

	// ensure the package path exists
	err = os.MkdirAll(espressoPath+"/cachedPackages", 0755)
	if err != nil {
		return err
	}

	// download the file
	err = util.DownloadFile(pkgPath, rdep.PackageVersion.ArtifactUrl)
	if err != nil {
		return err
	}
	return nil
}
