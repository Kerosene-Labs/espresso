package registry

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"hlafaille.xyz/espresso/v0/core/project"
	"hlafaille.xyz/espresso/v0/core/util"
)

// Package is a high level abstraction on top of the raw filesystem based registry caching. Package represents
// a package within the registry, but contains all required runtime information in a convenient struct.
type Package struct {
	Group       string
	Name        string
	Description string
	Versions    []PackageVersionDeclaration
	Declaration PackageDeclaration
	Registry    project.Registry
}

// PackageVersionDeclaration is the file format of a package version within a cached registry package
type PackageVersionDeclaration struct {
	Number                string   `yaml:"number"`
	ArtifactUrl           string   `yaml:"artifactUrl"`
	TransientDependencies []string `yaml:"transientDependencies"`
	IsAnnotationProcessor bool     `yaml:"isAnnotationProcessor"`
}

// PackageDeclaration is the file format of a package declaration
type PackageDeclaration struct {
	Name        string
	Description string                      `yaml:"description"`
	Versions    []PackageVersionDeclaration `yaml:"versions"`
}

// UnmarshalPackageDeclaration unmarshals a package declaration from yaml text
func UnmarshalPackageDeclaration(content string) (*PackageDeclaration, error) {
	var regPkg PackageDeclaration
	err := yaml.Unmarshal([]byte(content), &regPkg)
	if err != nil {
		return nil, err
	}
	return &regPkg, nil
}

// MarshalPackageDeclaration marshals a package declartion to yaml text
func MarshalPackageDeclaration(pkgDecl *PackageDeclaration) (string, error) {
	out, err := yaml.Marshal(pkgDecl)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// InvalidateRegistry invalidates a particular registry
func InvalidateRegistryCache(reg *project.Registry) error {
	// get our home dir
	cachePath, err := GetRegistryCachePath(reg)
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
	cachePath, err := GetRegistryCachePath(reg)
	if err != nil {
		return err
	}

	// if the cache exists, error out
	doesExist, err := util.DoesPathExist(cachePath)
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
	doesDepsExist, err := util.DoesPathExist(cachePath + "/lookup/espresso-registry-main/packages")
	if err != nil {
		fmt.Printf("An error occurred while reading the registry's lookup directory: %s\n", err)
	}
	if !doesDepsExist {
		fmt.Println("An eror occurred: this registry's lookup appears invalid: no packages directory")
	}

	fmt.Println("Done")
	return nil
}

// GetRegistryPackageDeclarations parses all package declarations within the cache for a given registry
func GetRegistryPackages(reg project.Registry) ([]Package, error) {
	// get the package group paths
	pkgGrpPths, err := walkRegistryLookup(reg)
	if err != nil {
		return nil, err
	}

	// iterate over each package group path, walk the directory and read the files
	var pkgs []Package = []Package{}
	for _, pkgGroupPth := range pkgGrpPths {
		// walk this particular package group for package declarations
		pkgDeclPaths, err := walkPackageGroup(pkgGroupPth)
		if err != nil {
			return nil, err
		}

		// iterate over the package declaration paths, unmarshal them
		for _, pkgDeclPath := range pkgDeclPaths {
			// read the package declaration
			declContent, err := os.ReadFile(pkgDeclPath)
			if err != nil {
				return nil, err
			}

			// get the unmarshalled file
			unmarshalledDecl, err := UnmarshalPackageDeclaration(string(declContent))
			if err != nil {
				return nil, err
			}

			// set our name
			split := strings.Split(pkgDeclPath, "/")
			unmarshalledDecl.Name = strings.TrimSuffix(split[len(split)-1], ".yml")

			// create our high level package
			splitPkgGrp := strings.Split(pkgGroupPth, "/")
			pkg := Package{
				Group:       splitPkgGrp[len(splitPkgGrp)-1],
				Name:        unmarshalledDecl.Name,
				Description: unmarshalledDecl.Description,
				Versions:    unmarshalledDecl.Versions,
				Declaration: *unmarshalledDecl,
				Registry:    reg,
			}
			pkgs = append(pkgs, pkg)
		}
	}
	return pkgs, nil
}

// GenerateSignature generates a unique signature of a package and version. This can be used to uniquely
// reference a local copy of a packages across registries.
func CalculatePackageSignature(dep *Package, version *PackageVersionDeclaration) string {
	if dep == nil || version == nil {
		panic("dep or version was nil")
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", dep.Group, dep.Name, version.Number))))
}
