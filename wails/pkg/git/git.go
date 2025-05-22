package git

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	appconfig "github.com/renatonmag/version-ctrls/pkg/config"
)

// Repository wraps a go-git repository
type LocalRepository struct {
	repo           *git.Repository
	Path           string
	MainRepoBranch string
	AppConfig      *appconfig.ApplicationConfig
}

// NewRepository creates a new Repository instance
func NewRepository() (*LocalRepository, error) {
	// parentRepoPath := GetParentRepoPath(path)
	// if parentRepoPath == "" {
	// 	return nil, fmt.Errorf("no parent git repository found")
	// }
	// _, err := os.ReadFile(filepath.Join(parentRepoPath, ".git", "HEAD"))
	// if err != nil {
	// 	return nil, fmt.Errorf("error reading file %s: %w", path, err)
	// }

	appConfig, err := appconfig.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}
	return &LocalRepository{
		Path:           appConfig.Integration.Path,
		MainRepoBranch: "",
		AppConfig:      appConfig,
	}, nil
}

func (repo *LocalRepository) SetMainRepoBranch(branch string) {
	repo.MainRepoBranch = GetLongName(branch)
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
func (repo *LocalRepository) Open() error {
	repository, err := git.PlainOpen(repo.AppConfig.Integration.Path)
	if err != nil {
		return fmt.Errorf("failed to open git repository: %w", err)
	}
	repo.repo = repository
	return nil
}

// CreateBranch creates a new branch with the given name.
// Returns an error if the creation fails.
func (repo *LocalRepository) CreateBranch(branchName string, srcBranch string) error {
	// Get the HEAD reference
	headRef, err := repo.repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD reference: %w", err)
	}
	commitHash := headRef.Hash()

	// If a source branch is specified, use its commit instead of HEAD
	if srcBranch != "" {
		fromBranchRefName := plumbing.NewBranchReferenceName(srcBranch)
		fromBranchRef, err := repo.repo.Reference(fromBranchRefName, true)
		if err != nil {
			return fmt.Errorf("failed to get source branch reference: %w", err)
		}
		commitHash = fromBranchRef.Hash()
	}
	newBranchRefName := plumbing.NewBranchReferenceName(branchName)

	// Create a new reference pointing to the empty commit
	newBranchRef := plumbing.NewHashReference(newBranchRefName, commitHash)

	// Check if the branch already exists
	exists := repo.BranchExists(branchName)
	if exists {
		return fmt.Errorf("branch '%s' already exists", branchName)
	}

	// Save the new branch reference to the repository's storage
	if err := repo.repo.Storer.SetReference(newBranchRef); err != nil {
		return fmt.Errorf("failed to create new branch: %w", err)
	}

	return nil
}

// DeleteBranch deletes a branch with the given name.
// Returns an error if the branch doesn't exist or if the deletion fails.
func (repo *LocalRepository) DeleteBranch(branchName string) error {
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

// RenameBranch renames a branch from oldName to newName.
// Returns an error if the source branch doesn't exist, the target branch already exists, or the operation fails.
func (repo *LocalRepository) RenameBranch(oldName, newName string) error {
	err := repo.CheckoutBranch("master")
	if err != nil {
		return fmt.Errorf("failed to checkout branch '%s': %w", oldName, err)
	}

	oldRefName := plumbing.NewBranchReferenceName(oldName)
	newRefName := plumbing.NewBranchReferenceName(newName)

	// Check if the old branch exists
	oldRef, err := repo.repo.Reference(oldRefName, false)
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			return fmt.Errorf("branch '%s' does not exist", oldName)
		}
		return fmt.Errorf("failed to get old branch reference: %w", err)
	}

	// Check if the new branch already exists
	_, err = repo.repo.Reference(newRefName, false)
	if err == nil {
		return fmt.Errorf("branch '%s' already exists", newName)
	} else if err != plumbing.ErrReferenceNotFound {
		return fmt.Errorf("failed to check new branch reference: %w", err)
	}

	// Create the new branch reference
	newRef := plumbing.NewHashReference(newRefName, oldRef.Hash())
	if err := repo.repo.Storer.SetReference(newRef); err != nil {
		return fmt.Errorf("failed to create new branch reference: %w", err)
	}

	// Remove the old branch reference
	if err := repo.repo.Storer.RemoveReference(oldRefName); err != nil {
		return fmt.Errorf("failed to remove old branch reference: %w", err)
	}

	return nil
}

func (repo *LocalRepository) CheckoutBranch(branchRef string) error {
	branchRefName := plumbing.NewBranchReferenceName(branchRef)
	wt, err := repo.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: branchRefName,
		Force:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to checkout branch '%s': %w", branchRef, err)
	}
	return nil
}

// CommitToBranch creates a commit on the specified branch without checking it out.
// It adds the specified file to the commit.
// Returns the commit hash as a string or an error if the operation fails.
func (repo *LocalRepository) CommitToBranch(branchName, filePath, commitMessage string, empty bool) (string, error) {
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
		Parents:           []plumbing.Hash{branchRef.Hash()},
		AllowEmptyCommits: empty,
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

// CommitToBranch creates a commit on the specified branch without checking it out.
// It adds the specified file to the commit.
// Returns the commit hash as a string or an error if the operation fails.
func (repo *LocalRepository) CommitAll(branchName, commitMessage string) (string, error) {
	// Check if the branch is currently checked out
	headRef, err := repo.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Get the current branch name
	currentBranch := headRef.Name().Short()

	// Check if we're trying to commit to the currently checked out branch
	fmt.Println("headRef", headRef)
	fmt.Println("branchName", branchName)
	fmt.Println("currentBranch", currentBranch)
	if currentBranch != branchName {
		return "", fmt.Errorf("cannot commit to branch '%s' because branch '%s' is currently checked out", branchName, currentBranch)
	}
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

	// Add all changes to the staging area
	if err := worktree.AddWithOptions(&git.AddOptions{
		All: true,
	}); err != nil {
		return "", fmt.Errorf("failed to add all changes to staging area: %w", err)
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

	return commitHash.String(), nil
}

// BranchExists checks if a branch with the given name exists in the repository.
// Returns true if the branch exists, false otherwise.
func (repo *LocalRepository) BranchExists(branchName string) bool {
	branchRefName := plumbing.NewBranchReferenceName(branchName)
	_, err := repo.repo.Reference(branchRefName, false)
	return err == nil
}

func (repo *LocalRepository) Push(remoteName string) error {
	if repo.repo == nil {
		return fmt.Errorf("repository is not open")
	}
	err := repo.repo.Push(&git.PushOptions{
		RemoteName: remoteName,
		Auth: &http.BasicAuth{
			Username: repo.AppConfig.Credentials.Username,
			Password: repo.AppConfig.Credentials.Password,
		},
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			return nil // Not an error, just nothing to push
		}
		return fmt.Errorf("failed to push to remote '%s': %w", remoteName, err)
	}
	return nil
}

// AddRemote adds a new remote to the repository with the given name and URL.
// Returns an error if the remote already exists or if the operation fails.
func (repo *LocalRepository) AddRemote(name, url string) error {
	if repo.repo == nil {
		return fmt.Errorf("repository is not open")
	}
	_, err := repo.repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	})
	if err != nil {
		return fmt.Errorf("failed to add remote '%s': %w", name, err)
	}
	return nil
}

// UpdateRemote updates the URL of an existing remote.
func (repo *LocalRepository) UpdateRemote(name, newURL string) error {
	if repo.repo == nil {
		return fmt.Errorf("repository is not open")
	}
	err := repo.repo.DeleteRemote(name)
	if err != nil {
		return fmt.Errorf("failed to delete existing remote '%s': %w", name, err)
	}
	_, err = repo.repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{newURL},
	})
	if err != nil {
		return fmt.Errorf("failed to add remote '%s' with new URL: %w", name, err)
	}
	return nil
}

// RemoveRemote removes a remote from the repository.
func (repo *LocalRepository) RemoveRemote(name string) error {
	if repo.repo == nil {
		return fmt.Errorf("repository is not open")
	}
	err := repo.repo.DeleteRemote(name)
	if err != nil {
		return fmt.Errorf("failed to remove remote '%s': %w", name, err)
	}
	return nil
}
