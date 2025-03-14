package git

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Repository wraps a go-git repository
type Repository struct {
	repo *git.Repository
	path string
}

// NewRepository creates a new Repository instance
func NewRepository(path string) *Repository {
	return &Repository{
		path: path,
	}
}

// InitRepository initializes a new git repository at the specified path.
// Returns an error if the initialization fails.
func InitRepository(path string) error {
	_, err := git.PlainInit(path, false)
	if err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}
	return nil
}

// RemoveDirectory removes a git repository at the specified path.
// This deletes the entire directory including the .git folder and all files.
// Returns an error if the deletion fails.
func RemoveDirectory(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("failed to remove directory: %w", err)
	}
	return nil
}

// Open opens a git repository at the specified path.
// Returns an error if the opening fails.
func (repo *Repository) Open() error {
	repository, err := git.PlainOpen(repo.path)
	if err != nil {
		return fmt.Errorf("failed to open git repository: %w", err)
	}
	repo.repo = repository
	return nil
}

// CreateBranch creates a new branch with the given name.
// Returns an error if the creation fails.
func (repo *Repository) CreateBranch(branchName string) error {
	newBranchRefName := plumbing.NewBranchReferenceName(branchName)

	// Get the worktree
	worktree, err := repo.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Create an empty commit
	commitHash, err := worktree.Commit("Initial empty commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "AutoCommit",
			Email: "autocommit@example.com",
			When:  time.Now(),
		},
		AllowEmptyCommits: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create an empty commit: %w", err)
	}

	// Create a new reference pointing to the empty commit
	newBranchRef := plumbing.NewHashReference(newBranchRefName, commitHash)

	// Check if the branch already exists
	_, err = repo.repo.Reference(newBranchRefName, false)
	if err == nil {
		return fmt.Errorf("branch '%s' already exists", branchName)
	} else if err != plumbing.ErrReferenceNotFound {
		return fmt.Errorf("failed to check if branch exists: %w", err)
	}

	// Save the new branch reference to the repository's storage
	if err := repo.repo.Storer.SetReference(newBranchRef); err != nil {
		return fmt.Errorf("failed to create new branch: %w", err)
	}

	return nil
}

// DeleteBranch deletes a branch with the given name.
// Returns an error if the branch doesn't exist or if the deletion fails.
func (repo *Repository) DeleteBranch(branchName string) error {
	branchRefName := plumbing.NewBranchReferenceName(branchName)

	// Check if the branch exists
	_, err := repo.repo.Reference(branchRefName, false)
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			return fmt.Errorf("branch '%s' does not exist", branchName)
		}
		return fmt.Errorf("failed to check if branch exists: %w", err)
	}

	// Delete the branch reference
	err = repo.repo.Storer.RemoveReference(branchRefName)
	if err != nil {
		return fmt.Errorf("failed to delete branch '%s': %w", branchName, err)
	}

	return nil
}

// CommitToBranch creates a commit on the specified branch without checking it out.
// It adds the specified file to the commit.
// Returns the commit hash as a string or an error if the operation fails.
func (repo *Repository) CommitToBranch(branchName, filePath, commitMessage string) (string, error) {
	// Get the branch reference
	branchRefName := plumbing.NewBranchReferenceName(branchName)
	branchRef, err := repo.repo.Reference(branchRefName, true)
	if err != nil {
		return "", fmt.Errorf("failed to get branch reference: %w", err)
	}

	// Get the worktree
	worktree, err := repo.repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	// Add the specified file to the staging area
	if err := worktree.AddWithOptions(&git.AddOptions{
		Path: filePath,
	}); err != nil {
		return "", fmt.Errorf("failed to add file '%s' to staging area: %w", filePath, err)
	}

	// Create the commit
	commitOptions := &git.CommitOptions{
		Author: &object.Signature{
			Name:  "AutoCommit",
			Email: "autocommit@example.com",
			When:  time.Now(),
		},
		Parents: []plumbing.Hash{branchRef.Hash()},
	}

	// Commit to the branch
	commitHash, err := worktree.Commit(commitMessage, commitOptions)
	if err != nil {
		return "", fmt.Errorf("failed to commit to branch '%s': %w", branchName, err)
	}

	// Update the branch reference to point to the new commit
	newRef := plumbing.NewHashReference(branchRefName, commitHash)
	if err := repo.repo.Storer.SetReference(newRef); err != nil {
		return "", fmt.Errorf("failed to update branch reference: %w", err)
	}

	return commitHash.String(), nil
}
