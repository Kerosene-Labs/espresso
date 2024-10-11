package project

import "fmt"

// GetVersionAsString gets a project version as a string.
func GetVersionAsString(version Version) string {
	versionString := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)
	if version.Hotfix != nil {
		versionString += *version.Hotfix
	}
	return versionString
}
