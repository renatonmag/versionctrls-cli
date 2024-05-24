package repository

import (
	"errors"
	"path/filepath"
)

// GetRepoRoot checks if the current directory is the root of a Git repository
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

	return "", errors.New("you are not in a git repository root")
}

// ConcatenateSubmodulePath concatenates the submodule path with the repository root
func (r *Repository) IntegrationSubmodulePath() (string, error) {
	repoRoot, err := r.GetRepoRoot()
	if err != nil {
		return "", err
	}

	// Concatenate the submodule path with the repository root
	fullSubmodulePath := filepath.Join(repoRoot, r.submodulePath)
	return fullSubmodulePath, nil
}
