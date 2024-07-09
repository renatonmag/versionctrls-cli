package repository

import (
	"github.com/go-git/go-git/v5/plumbing"
)

// branchExists checks if a branch with the given name exists in the repository.
func (r *Repository) BranchExists(branchName string) (bool, error) {
	// Get the reference name for the branch.
	refName := plumbing.NewBranchReferenceName(branchName)

	// Try to get the reference.
	_, err := r.repo.Reference(refName, false)
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
