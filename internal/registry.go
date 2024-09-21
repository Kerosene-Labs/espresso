package internal

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GetCachePath gets the full cache path from the registry
func GetCachePath(reg *Registry) (string, error) {
	// get our home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.espresso/registries/" + reg.Name, nil
}

// InvalidateRegistry invalidates a particular registry
func InvalidateRegistryCache(reg *Registry) error {
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
func CacheRegistry(reg *Registry) error {
	// get our cache path
	cachePath, err := GetCachePath(reg)
	if err != nil {
		return err
	}

	// if the cache exists, error out
	doesExist, err := DoesFileExist(cachePath)
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
	err = DownloadFile(cachePath+"/archive.zip", reg.Url)
	if err != nil {
		return err
	}

	// extract the archive
	fmt.Println("Extracting")
	unzip(cachePath+"/archive.zip", cachePath+"/lookup")

	// check if the registry lookup contains a dependencies folder
	doesDepsExist, err := DoesFileExist(cachePath + "/lookup/espresso-registry-main/dependencies")
	if err != nil {
		fmt.Printf("An error occurred while reading the registry's lookup directory: %s\n", err)
	}
	if !doesDepsExist {
		fmt.Println("An eror occurred: this registry's lookup appears invalid: no dependencies directory")
	}

	fmt.Println("Done")
	return nil
}
