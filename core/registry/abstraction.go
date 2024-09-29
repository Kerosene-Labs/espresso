package registry

// Package is a high level abstraction on top of the raw filesystem based registry caching. Package represents
// a package within the registry, but contains all required runtime information in a convenient struct.
type Package struct {
	Group       string
	Name        string
	Description string
	Versions    []PackageVersionDeclaration
	Delcaration PackageDeclaration
}
