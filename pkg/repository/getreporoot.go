package repository

import (
	"errors"
	"path/filepath"
)

// IsGitRepoRoot checks if the current directory is the root of a Git repository
func (r *Repository) GetRepoRoot() (string, error) {
	// Get the repository worktree
	worktree, err := r.repo.Worktree()
	if err != nil {
		return "", err
	}

	// Get the absolute path of the current directory
	currentDir, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}

	// Get the absolute path of the repository root
	repoRoot := worktree.Filesystem.Root()
	repoRoot, err = filepath.Abs(repoRoot)
	if err != nil {
		return "", err
	}

	if currentDir == repoRoot {
		return repoRoot, nil
	}

	return "", errors.New("You are not in a Git repository root")
}
