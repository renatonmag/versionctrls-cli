package git

import (
	"os"
	"path/filepath"
	"strings"
)

func GetShortName(ref_str string) string {
	ref := strings.TrimSpace(ref_str)
	if strings.HasPrefix(ref, "refs/heads/") {
		return strings.TrimPrefix(ref, "refs/heads/")
	} else if strings.HasPrefix(ref, "ref: refs/heads/") {
		return strings.TrimPrefix(ref, "ref: refs/heads/")
	}
	return ""
}

func GetLongName(ref_str string) string {
	ref := strings.TrimSpace(ref_str)
	if strings.HasPrefix(ref, "ref: ") {
		return strings.TrimPrefix(ref, "ref: ")
	}
	return ""
}

// GetParentRepoPath returns the path to the nearest parent git repository
// Returns an empty string if no parent git repository is found
func GetParentRepoPath(path string) string {
	currentPath := path
	for currentPath != "/" && currentPath != "" {
		currentPath = filepath.Dir(currentPath)
		if _, err := os.Stat(filepath.Join(currentPath, ".git")); err == nil {
			return currentPath
		}
	}
	return ""
}
