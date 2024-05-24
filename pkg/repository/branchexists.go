package repository

import "github.com/go-git/go-git/v5/plumbing"

// BranchExists checks if a branch already exists in the repository
func (r *Repository) BranchExists(branchName string) (bool, error) {
	_, err := r.repo.Reference(plumbing.NewBranchReferenceName(branchName), true)
	if err == plumbing.ErrReferenceNotFound {
		return false, nil
	}
	return err == nil, err
}
