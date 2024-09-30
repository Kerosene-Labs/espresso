package dependency

import (
	"fmt"

	"hlafaille.xyz/espresso/v0/core/project"
	"hlafaille.xyz/espresso/v0/core/registry"
	"hlafaille.xyz/espresso/v0/core/util"
)

// FetchDependency fetches a version of a given package.
func FetchDependency(pkg *registry.Package, version *registry.PackageVersionDeclaration) error {
	if pkg == nil || version == nil {
		panic("pkg or version was nil")
	}
	util.DownloadFile()
	return nil
}

// GetVersionAsString gets a project version as a string.
func GetVersionAsString(version *project.Version) string {
	if version == nil {
		panic("version is nil")
	}
	return fmt.Sprintf("%s.%s.%s%s", version.Major, version.Minor, version.Patch, version.Hotfix)
}

func GenerateSignature(dep *project.Dependency)
