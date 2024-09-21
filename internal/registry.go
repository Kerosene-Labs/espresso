package internal

import (
	"errors"
	"os"
	"strings"
)

// GetRegistryNameFromUrl gets the registry name from the URL. This is achieved by splitting the URL by "/"
// and returning the last element in the slice. For example, https://github.com/Kerosene-Labs/espresso-registry
// would have "espress-registry" as its name.
func GetRegistryNameFromUrl(reg *Registry) string {
	split := strings.Split(reg.Url, "/")
	return split[len(split)-1]
}

// GetCachePath gets the full cache path from the registry
func GetCachePath(reg *Registry) (string, error) {
	// get our home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.espresso/registries/" + GetRegistryNameFromUrl(reg), nil
}

// InvalidateRegistry invalidates a particular registry
func InvalidateRegistry(reg *Registry) error {
	// get our home dir
	cachePath, err := GetCachePath(reg)
	if err != nil {
		return err
	}

	// delete the registry cache
	os.RemoveAll(cachePath)
	return nil
}

// Check if this registry is already cached
func CacheRegistry(reg *Registry) error {
	// get our cache path
	cachePath, err := GetCachePath(reg)
	if err != nil {
		return err
	}

	// if the cache exists, error out
	doesExist, err := DoesFileExist(cachePath)
	if doesExist {
		return errors.New("cache exists: must be invalidated or not exist")
	}

	// TODO: implement non-github solution
	// ensure this registry starts with github
	if !strings.HasPrefix(reg.Url, "https://github.com") {
		return errors.New("unsupported registry: for now, registries must be public github repositories. this will change in a later release.")
	}

	// call github to get this repo's zip
	// https://github.com/OWNER/REPO/archive/refs/heads/BRANCH.zip

	return nil
}
