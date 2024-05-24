package repository

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Commit changes to the repository
func (r Repository) CommitChanges(path, branchName string) error {
	if r.repo == nil {
		return errors.New("no repository opened")
	}

	// Check if the branch already exists
	exists, err := r.BranchExists(branchName)
	if err != nil {
		log.Fatalf("could not check if branch exists: %v", err)
		return err
	}

	// Switch to the branch, or create it if it doesn't exist
	worktree, err := r.repo.Worktree()
	if err != nil {
		return err
	}

	if exists {
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branchName),
		})
		log.Printf("Switching to branch %s\n", branchName)
	} else {
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branchName),
			Create: true,
		})
	}
	if err != nil {
		log.Fatalf("could not switch to branch %s: %v", branchName, err)
		return err
	}

	// Add the file to the staging area
	_, err = worktree.Add(path)
	if err != nil {
		log.Fatalf("could not add file to staging area: %v", err)
		return err
	}

	authorName, authorEmail, err := r.GetGitUserInfo()
	if err != nil {
		return err
	}

	// Commit the changes
	commit, err := worktree.Commit(path+"-v0.1.0", &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("could not commit changes: %w", err)
	}

	// Print the commit hash
	obj, err := r.repo.CommitObject(commit)
	if err != nil {
		return fmt.Errorf("could not get commit object: %w", err)
	}
	fmt.Println("Commit successful:", obj.Hash)

	return nil
}
