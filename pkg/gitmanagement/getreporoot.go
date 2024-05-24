package gitmanagement

import (
	"errors"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

// IsGitRepoRoot checks if the current directory is the root of a Git repository
func GetRepoRoot() (string, error) {
	// Open the current directory as a Git repository
	repo, err := git.PlainOpen(".")
	if err != nil {
		// Not a Git repository
		return "", errors.New("You are not in a Git repository root")
	}

	// Get the repository worktree
	worktree, err := repo.Worktree()
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
