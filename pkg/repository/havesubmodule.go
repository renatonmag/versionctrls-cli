package repository

import (
	"errors"

	"github.com/go-git/go-git/v5"
)

// SubmoduleExists checks if a submodule with the given name already exists
func (r *Repository) SubmoduleExists(submoduleName string) (bool, error) {
	if r.repo == nil {
		return false, errors.New("no repository opened")
	}

	// Get the worktree
	worktree, err := r.repo.Worktree()
	if err != nil {
		return false, err
	}

	// Get the submodules
	submodules, err := worktree.Submodules()
	if err != nil {
		return false, err
	}

	// Check if the submodule with the given name exists
	for _, submodule := range submodules {
		if submodule.Config().Name == submoduleName {
			// Check if the submodule is initialized
			submoduleRepo, err := submodule.Repository()
			if err == git.ErrSubmoduleNotInitialized {
				return false, nil
			}
			return submoduleRepo != nil, nil
		}
	}

	return false, nil
}
