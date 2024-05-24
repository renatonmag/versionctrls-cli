package gitmanagement

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Commit changes to the repository
func commitChanges(repoPath, message, authorName, authorEmail string) error {
	// Open the existing repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("could not open repository: %w", err)
	}

	// Get the working directory for the repository
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("could not get worktree: %w", err)
	}

	// Add changes to the staging area
	if _, err = worktree.Add("."); err != nil {
		return fmt.Errorf("could not add changes: %w", err)
	}

	// Commit the changes
	commit, err := worktree.Commit(message, &git.CommitOptions{
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
	obj, err := repo.CommitObject(commit)
	if err != nil {
		return fmt.Errorf("could not get commit object: %w", err)
	}
	fmt.Println("Commit successful:", obj.Hash)

	return nil
}
