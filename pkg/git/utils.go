package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetShortName returns the short name of a git reference
// Returns an empty string if the reference is not a branch
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

func directoryExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		// File or directory exists
		fileInfo, err := os.Stat(path)
		if err != nil {
			// Handle the error of getting file info (unlikely, but possible)
			fmt.Println("Error getting file info:", err) // Log or handle appropriately
			return false                                 // Indicate failure, don't assume existence
		}
		return fileInfo.IsDir() // Check if it's a directory
	} else if os.IsNotExist(err) {
		// Directory does not exist
		return false
	} else {
		// Some other error occurred (e.g., permission denied)
		fmt.Println("Error checking directory:", err) // Log or handle appropriately
		return false                                  // Indicate failure, don't assume existence
	}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true // file exists
	} else if os.IsNotExist(err) {
		return false
	} else {
		fmt.Printf("error checking file existence: %w", err)
	}
	return false
}

func getBranchName(mainRepoBranch string, path string) string {
	return fmt.Sprintf("{%s}/{%s}", strings.ReplaceAll(GetShortName(mainRepoBranch), "/", "---"), strings.ReplaceAll(path, "/", "---"))
}

// func replaceHardLink(repo *Repository, event *fsbroker.FSEvent) {
// 	branchName := getBranchName(repo, event)

// 	exists := repo.BranchExists(branchName)
// 	if !exists {
// 		return fmt.Errorf("branch %s does not exist", branchName)
// 	}

// }
