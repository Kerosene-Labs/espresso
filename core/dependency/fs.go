// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package dependency

import (
	"os"

	"kerosenelabs.com/espresso/core/util"
)

// EnsurePackagesDirectory ensures the packages directory exists under the given registry cache.
func EnsurePackagesDirectory() error {
	ePath, err := util.GetEspressoDirectoryPath()
	if err != nil {
		return err
	}
	ePathExist, err := util.DoesPathExist(ePath)
	if err != nil {
		return err
	}
	if !ePathExist {
		os.MkdirAll(ePath, 0755)
	}
	return nil
}
