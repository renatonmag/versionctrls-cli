package repository

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CreateEmptyBranch creates a new branch in the given repository and makes an initial commit with an empty README.md file
func (r Repository) CreateEmptyBranch(branchName string) error {
	if r.repo == nil {
		return errors.New("no repository opened")
	}

	// Check if the branch already exists
	branchRef := plumbing.NewBranchReferenceName(branchName)
	_, err := r.repo.Reference(branchRef, true)
	if err == nil {
		// Branch already exists
		// log.Printf("Branch %s already exists\n", branchName)
		return nil
	}

	// Create the branch reference
	headRef, err := r.repo.Head()
	if err != nil {
		return err
	}

	// Create a new branch pointing to the current HEAD
	ref := plumbing.NewHashReference(branchRef, headRef.Hash())
	if err := r.repo.Storer.SetReference(ref); err != nil {
		return err
	}

	// Checkout the new branch
	worktree, err := r.repo.Worktree()
	if err != nil {
		return err
	}
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branchRef,
		Create: false,
	})
	if err != nil {
		return err
	}

	// Create and stage the README.md file
	readmePath := filepath.Join(worktree.Filesystem.Root(), "README.md")
	err = os.WriteFile(readmePath, []byte{}, 0644)
	if err != nil {
		return err
	}

	_, err = worktree.Add("README.md")
	if err != nil {
		return err
	}

	// Create an initial commit
	_, err = worktree.Commit("Initial commit with README.md", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Author Name",
			Email: "author@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
