package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// Folder represents a directory with a map of its children
type Folder struct {
	Path     string
	Name     string
	IsDir    bool
	Children []*Folder
}

// NewFolder creates a new folder
func NewFolder(path string, name string, isDir bool) *Folder {
	return &Folder{Path: path, Name: name, IsDir: isDir, Children: make([]*Folder, 0)}
}

// AddChild adds a child to a directory
func (f *Folder) AddChild(path string, name string, isDir bool) *Folder {
	child := NewFolder(path, name, isDir)
	f.Children = append(f.Children, child)
	return child
}

func CreateDirectoryStructure(path string) *Folder {
	root := NewFolder(path, filepath.Base(path), true)
	folderMap := map[string]*Folder{
		path: root,
	}

	err := filepath.Walk(path, func(currPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if currPath == path {
			return nil // skip root, already added
		}
		parentPath := filepath.Dir(currPath)
		parentFolder, ok := folderMap[parentPath]
		if !ok {
			// Optionally log or handle this error
			return fmt.Errorf("parent folder not found for path: %s", currPath)
		}
		child := parentFolder.AddChild(currPath, info.Name(), info.IsDir())
		if info.IsDir() {
			folderMap[currPath] = child
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", path, err)
		return nil
	}
	return root
}

// PrintTree prints the folder structure
func (f *Folder) PrintTree(indent string) {
	fmt.Println(indent + f.Name)
	for _, child := range f.Children {
		child.PrintTree(indent + "  ")
	}
}
