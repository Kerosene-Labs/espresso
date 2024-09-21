package dependency

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip vulnerability: https://snyk.io/research/zip-slip-vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create the file's directory if it doesn't exist
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Create the file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Copy the file's contents to the created file
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

// Check if this registry is already cached
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
	unzip(cachePath+"/archive.zip", cachePath+"/lookup")

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

// getDirectoriesInRegistryCache gets all directories (aka groupId's) within the specified registry cache
func getDirectoriesInRegistryCache(reg project.Registry) ([]string, error) {
	// get the cache path
	cachePath, err := GetCachePath(&reg)
	if err != nil {
		return []string{}, err
	}

	// walk the directory for all groupId's
	var dirs []string = []string{}
	err = filepath.Walk(cachePath+"/lookup/espresso-registry-main/dependencies", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	return dirs, nil
}

// GetRegistryPackages parses all packages within the cache for a given registry
func GetRegistryPackages(reg project.Registry) ([]RegistryPackage, error) {
	// get the directories
	dirs, err := getDirectoriesInRegistryCache(reg)
	if err != nil {
		return []RegistryPackage{}, err
	}

	fmt.Printf("%s", dirs)
	return []RegistryPackage{}, nil
}

// ResolveDependency determines which registry is the best fit for a given dependency
// func ResolveDependency(dep Dependency, regs []Registry) (Registry, error) {
// 	for _, reg := range regs {

// 	}
// }
