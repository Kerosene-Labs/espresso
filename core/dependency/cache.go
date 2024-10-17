// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package dependency

import (
	"os"

	"kerosenelabs.com/espresso/core/util"
)

// CacheResolvedDependency fetches the resolved dependency from the internet
func CacheResolvedDependency(resolvedDependency ResolvedDependency) error {
	// get where we should store this package
	espressoPath, err := util.GetEspressoDirectoryPath()
	if err != nil {
		return err
	}

	// ensure the package path exists
	err = os.MkdirAll(espressoPath+"/cachedPackages", 0755)
	if err != nil {
		return err
	}

	// get our package cache path
	pkgPath, err := resolvedDependency.GetCachePath()
	if err != nil {
		return err
	}

	// download the file
	err = util.DownloadFile(pkgPath.Absolute, resolvedDependency.PackageVersion.ArtifactUrl)
	if err != nil {
		return err
	}
	return nil
}
