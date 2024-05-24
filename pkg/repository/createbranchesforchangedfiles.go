package repository

import (
	"path/filepath"
	"strings"
)

// CreateEmptyBranchesForChangedFiles creates an empty branch for each changed file with its path as the branch name
func (r Repository) CreateEmptyBranchesForChangedFiles() error {
	changedFiles, err := r.GetChangedFiles()
	if err != nil {
		return err
	}

	for _, file := range changedFiles {
		// Replace path separators with a character allowed in branch names (e.g., "-")
		branchName := strings.ReplaceAll(file, string(filepath.Separator), "-")

		err := r.CreateEmptyBranch(branchName)
		if err != nil {
			return err
		}
	}

	return nil
}
