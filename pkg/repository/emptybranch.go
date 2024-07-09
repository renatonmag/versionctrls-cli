package repository

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/filesystem/dotgit"
)

// CreateEmptyBranch creates a new branch in the given repository and makes an initial commit with an empty README.md file
func (r Repository) CreateEmptyBranch(branchName string) error {
	if r.repo == nil {
		return errors.New("no repository opened")
	}

	// // Check if the branch already exists
	exists, _ := r.BranchExists(branchName)
	if exists {
		// Branch already exists
		log.Printf("Branch %s already exists\n", branchName)
		return nil
	}

	refName := plumbing.NewBranchReferenceName(branchName)
	newBranchRef := plumbing.NewHashReference(refName, plumbing.ZeroHash)

	// Set the new branch reference in the repository.
	err := r.repo.Storer.SetReference(newBranchRef)
	if err != nil {
		log.Fatalf("Failed to create branch reference: %s", err)
	}

	// Check out the new branch.
	worktree, err := r.repo.Worktree()
	if err != nil {
		log.Fatalf("Failed to get worktree: %s", err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: refName,
		Create: false,
		Force:  false,
	})
	if err != nil {
		log.Fatalf("Failed to checkout new branch: %s", err)
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

	name, email, _ := r.GetGitUserInfo()
	// Create an initial commit
	commit, err := worktree.Commit("Initial commit with README.md", &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	fmt.Printf("Initial commit created: %s\n", commit)

	return nil
}

func (r Repository) CreateEmptyBranch2(branchName string) error {
	if r.repo == nil {
		return errors.New("no repository opened")
	}

	// // Check if the branch already exists
	exists, _ := r.BranchExists(branchName)
	if exists {
		// Branch already exists
		log.Printf("Branch %s already exists\n", branchName)
		return nil
	}

	// Create a new branch reference.
	refName := plumbing.NewBranchReferenceName(branchName)

	// Get the worktree.
	worktree, err := r.repo.Worktree()
	if err != nil {
		log.Fatalf("Failed to get worktree: %s", err)
	}

	// Checkout to the new branch.
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: refName,
		Create: true,
		Force:  true, // Force checkout to allow empty worktree
	})
	if err != nil {
		log.Fatalf("Failed to checkout new branchhhh: %s", err)
	}

	repoPath, _ := r.GetRepoRoot()
	gitDir := filepath.Join(repoPath, ".git")
	fs := osfs.New(gitDir)
	dotGit := dotgit.New(fs)

	// Step 3: Initialize cache.Object
	objectCache := cache.NewObjectLRUDefault()

	// Step 4: Call NewObjectStorage function
	objectStorage := filesystem.NewObjectStorage(dotGit, objectCache)

	// Create a new tree object
	tree := &object.Tree{}

	name, email, _ := r.GetGitUserInfo()
	// Create a new commit object
	commit := &object.Commit{
		Author: object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
		Committer: object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
		Message:  "This is a dangling commit",
		TreeHash: tree.Hash,
	}
	encObject := objectStorage.NewEncodedObject()
	err = commit.Encode(encObject)
	if err != nil {
		log.Fatalf("Failed to encode the commit object: %v", err)
	}

	// Encode the commit object to get the object ID (hash)
	objID, err := r.repo.Storer.SetEncodedObject(encObject)
	if err != nil {
		log.Fatalf("Failed to store the commit object: %v", err)
	}

	fmt.Printf("Created dangling commit with hash: %s\n", objID.String())

	// Create and stage the README.md file
	readmePath := filepath.Join(worktree.Filesystem.Root(), "README.md")
	err = os.WriteFile(readmePath, []byte(fmt.Sprintf("# %s", branchName)), 0644)
	if err != nil {
		log.Fatalf("Failed to write new README.md: %s", err)
		return err
	}

	// Add the new file to the staging area.
	_, err = worktree.Add("README.md")
	if err != nil {
		log.Fatalf("Failed to add file to the staging area: %s", err)
	}

	// name, email, _ := r.GetGitUserInfo()
	// Create an initial commit.
	second_commit, err := worktree.Commit("Add README.md", &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
		Parents: []plumbing.Hash{objID},
	})
	if err != nil {
		log.Fatalf("Failed to create initial commit: %s", err)
		return err
	}

	// Print the commit hash.
	fmt.Printf("Initial commit created: %s\n", second_commit)

	mainBranch := plumbing.NewBranchReferenceName("main")

	_ = worktree.Checkout(&git.CheckoutOptions{
		Branch: mainBranch,
		Create: false,
		Force:  false, // Force checkout to allow empty worktree
	})

	return nil
}

// // CreateEmptyBranch creates a new branch in the given repository and makes an initial commit with an empty README.md file
// func (r Repository) CreateEmptyBranch(branchName string) error {
// 	if r.repo == nil {
// 		return errors.New("no repository opened")
// 	}

// 	// Check if the branch already exists
// 	branchRef := plumbing.NewBranchReferenceName(branchName)
// 	_, err := r.repo.Reference(branchRef, true)
// 	if err == nil {
// 		// Branch already exists
// 		log.Printf("Branch %s already exists\n", branchName)
// 		return nil
// 	}

// 	// Create the branch reference
// 	headRef, err := r.repo.Head()
// 	if err != nil {
// 		return err
// 	}

// 	// Create a new branch pointing to the current HEAD
// 	ref := plumbing.NewHashReference(branchRef, headRef.Hash())
// 	if err := r.repo.Storer.SetReference(ref); err != nil {
// 		return err
// 	}

// 	// Checkout the new branch
// 	worktree, err := r.repo.Worktree()
// 	if err != nil {
// 		return err
// 	}
// 	err = worktree.Checkout(&git.CheckoutOptions{
// 		Branch: branchRef,
// 		Create: false,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	// Create and stage the README.md file
// 	readmePath := filepath.Join(worktree.Filesystem.Root(), "README.md")
// 	err = os.WriteFile(readmePath, []byte{}, 0644)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = worktree.Add("README.md")
// 	if err != nil {
// 		return err
// 	}

// 	// Create an initial commit
// 	_, err = worktree.Commit("Initial commit with README.md", &git.CommitOptions{
// 		Author: &object.Signature{
// 			Name:  "Author Name",
// 			Email: "author@example.com",
// 			When:  time.Now(),
// 		},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
