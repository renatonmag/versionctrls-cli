package repository

import (
	"errors"

	"github.com/go-git/go-git/v5"
)

// GetChangedFiles returns a list of changed files in the repository
func (r Repository) GetChangedFiles() ([]string, error) {
	var changedFiles []string
	if r.repo == nil {
		return []string{}, errors.New("no repository opened")
	}

	worktree, err := r.repo.Worktree()
	if err != nil {
		return []string{}, err
	}

	status, err := worktree.Status()
	if err != nil {
		return []string{}, err
	}

	for file, fileStatus := range status {
		if fileStatus.Worktree != git.Unmodified {
			changedFiles = append(changedFiles, file)
		}
	}

	return changedFiles, nil
}
