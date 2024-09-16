package project

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// SourceFIle represents a .java source file on the filesystem
type SourceFile struct {
	Path    string
	Content string
}

// readSourceFile reads the file at the given path, or returns an error
func readSourceFile(path string) (*SourceFile, error) {
	return nil, nil
}

// IntrospectSourceFiles iterates over the project's base package looking for .java files
// TODO: make this a goroutine?
func DiscoverSourceFiles(cfg *ProjectConfig) ([]SourceFile, error) {
	// get the source path
	srcPath, err := GetSourcePath(cfg)
	if err != nil {
		return nil, err
	}

	// iterate recursively over child directories
	var files []SourceFile = []SourceFile{}
	err = filepath.Walk(*srcPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".java") {
			text, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			files = append(files, SourceFile{
				Path:    path,
				Content: string(text),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
