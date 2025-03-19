package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// Folder represents a directory with a map of its children
type Folder struct {
	Name     string
	IsDir    bool
	Children map[string]*Folder
}

// NewFolder creates a new folder
func NewFolder(name string, isDir bool) *Folder {
	return &Folder{Name: name, IsDir: isDir, Children: make(map[string]*Folder)}
}

// AddChild adds a child to a directory
func (f *Folder) AddChild(name string, isDir bool) *Folder {
	if !f.IsDir {
		return nil
	}
	child := NewFolder(name, isDir)
	f.Children[name] = child
	return child
}

func CreateDirectoryStructure(path string) *Folder {
	root := NewFolder(path, true)
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return nil
	})
	return root
}

// PrintTree prints the folder structure
func (f *Folder) PrintTree(indent string) {
	fmt.Println(indent + f.Name)
	for _, child := range f.Children {
		child.PrintTree(indent + "  ")
	}
}
