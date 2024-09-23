package dependency

import (
	"errors"
	"fmt"
	"os"

	"hlafaille.xyz/espresso/v0/core/project"
	"hlafaille.xyz/espresso/v0/core/util"
)

type RegistryPackageVersion struct {
	Number                string   `yaml:"number"`
	ArtifactUrl           string   `yaml:"artifactUrl"`
	TransientDependencies []string `yaml:"transientDependencies"`
	IsAnnotationProcessor bool     `yaml:"isAnnotationProcessor"`
}

type RegistryPackage struct {
	Name        string                   `yaml:"name"`
	Description string                   `yaml:"description"`
	Versions    []RegistryPackageVersion `yaml:"versions"`
}

// GetCachePath gets the full cache path from the registry
func GetCachePath(reg *project.Registry) (string, error) {
	// get our home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.espresso/registries/" + reg.Name, nil
}

func GetCacheDependenciesPath(reg *project.Registry) (string, error) {
	// get our home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.espresso/registries/" + reg.Name + "espresso-registry-main/dependencies", nil
}

// InvalidateRegistry invalidates a particular registry
func InvalidateRegistryCache(reg *project.Registry) error {
	// get our home dir
	cachePath, err := GetCachePath(reg)
	if err != nil {
		return err
	}

	// delete the registry cache
	os.RemoveAll(cachePath)
	return nil
}

// CacheRegistry downloads a zip archive representing an espresso registry and extracts it to the proper directory
func CacheRegistry(reg *project.Registry) error {
	// get our cache path
	cachePath, err := GetCachePath(reg)
	if err != nil {
		return err
	}

	// if the cache exists, error out
	doesExist, err := util.DoesFileExist(cachePath)
	if err != nil {
		return err
	}
	if doesExist {
		return errors.New("cache exists: must be invalidated or not exist")
	}

	// create the registry's directory
	err = os.MkdirAll(cachePath, 0755)
	if err != nil {
		return err
	}

	// download the registry archive
	err = util.DownloadFile(cachePath+"/archive.zip", reg.Url)
	if err != nil {
		return err
	}

	// extract the archive
	fmt.Println("Extracting")
	util.Unzip(cachePath+"/archive.zip", cachePath+"/lookup")

	// check if the registry lookup contains a dependencies folder
	doesDepsExist, err := util.DoesFileExist(cachePath + "/lookup/espresso-registry-main/dependencies")
	if err != nil {
		fmt.Printf("An error occurred while reading the registry's lookup directory: %s\n", err)
	}
	if !doesDepsExist {
		fmt.Println("An eror occurred: this registry's lookup appears invalid: no dependencies directory")
	}

	fmt.Println("Done")
	return nil
}

// GetRegistryPackages parses all packages within the cache for a given registry
func GetRegistryPackages(reg project.Registry) ([]RegistryPackage, error) {
	// get the directories
	dirs, err := getPackageGroupPathsInRegCache(reg)
	if err != nil {
		return []RegistryPackage{}, err
	}

	fmt.Printf("%s\n", dirs)
	return []RegistryPackage{}, nil
}

// QueryRegistryCache queries the registry for packages matching the search term
func QueryRegistry(reg *project.Registry) (string, error) {
	if reg == nil {
		panic("programming error: reg was a nil pointer")
	}
	return "", nil

	//
}
