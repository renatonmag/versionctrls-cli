package repository

import (
	"log"
	"path/filepath"
	"strings"
)

// CreateCommitForChangedFiles creates a commit for each changed file in its own branch
func (r *Repository) CommitChangedFiles() error {
	changedFiles, err := r.GetChangedFiles()
	if err != nil {
		log.Printf("Error getting changed files: %s\n", err)
		return err
	}

	for _, file := range changedFiles {
		// Replace path separators with a character allowed in branch names (e.g., "-")
		branchName := strings.ReplaceAll(file, string(filepath.Separator), "-")

		err = r.CommitChanges(file, branchName)
		if err != nil {
			return err
		}
	}

	return nil
}
