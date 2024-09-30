package dependency

import (
	"os"

	"hlafaille.xyz/espresso/v0/core/util"
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
