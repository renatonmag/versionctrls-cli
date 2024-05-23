// File: main.go
// package main

// import (
//     "fmt"
//     "log"
//     "os"
//     "time"

//     "github.com/go-git/go-git/v5"
//     "github.com/go-git/go-git/v5/config"
//     "github.com/go-git/go-git/v5/plumbing/object"
// )

// // Commit changes to the repository
// func commitChanges(repoPath, message, authorName, authorEmail string) error {
//     // Open the existing repository
//     repo, err := git.PlainOpen(repoPath)
//     if err != nil {
//         return fmt.Errorf("could not open repository: %w", err)
//     }

//     // Get the working directory for the repository
//     worktree, err := repo.Worktree()
//     if err != nil {
//         return fmt.Errorf("could not get worktree: %w", err)
//     }

//     // Add changes to the staging area
//     if _, err = worktree.Add("."); err != nil {
//         return fmt.Errorf("could not add changes: %w", err)
//     }

//     // Commit the changes
//     commit, err := worktree.Commit(message, &git.CommitOptions{
//         Author: &object.Signature{
//             Name:  authorName,
//             Email: authorEmail,
//             When:  time.Now(),
//         },
//     })
//     if err != nil {
//         return fmt.Errorf("could not commit changes: %w", err)
//     }

//     // Print the commit hash
//     obj, err := repo.CommitObject(commit)
//     if err != nil {
//         return fmt.Errorf("could not get commit object: %w", err)
//     }
//     fmt.Println("Commit successful:", obj.Hash)

//     return nil
// }

// func main() {
//     // Replace these variables with your own values
//     repoPath := "/path/to/your/repo"
//     commitMessage := "Your commit message"
//     authorName := "Your Name"
//     authorEmail := "your-email@example.com"

//     // Commit changes
//     if err := commitChanges(repoPath, commitMessage, authorName, authorEmail); err != nil {
//         log.Fatalf("Failed to commit changes: %v", err)
//     }
// }

// File: main.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// RunGitCommand runs a Git command and returns the output
func RunGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run git %v: %v\nOutput: %s", args, err, string(output))
	}
	return string(output), nil
}

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

// // Clone a repository using go-git
func cloneRepoWithGoGit(url, path, username, password string) {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Repository cloned successfully using go-git.")
}

// Add a submodule using the Git CLI
func addSubmoduleWithGitCLI(repoPath, submodulePath, submoduleURL string) {
	// Change directory to the cloned repository
	if err := os.Chdir(repoPath); err != nil {
		log.Fatalf("Failed to change directory: %v", err)
	}

	// Add a submodule
	fmt.Println("Adding submodule using Git CLI...")
	if output, err := RunGitCommand("submodule", "add", submoduleURL, submodulePath); err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, output)
	} else {
		fmt.Println("Submodule added successfully using Git CLI.")
	}

	// Initialize submodules
	fmt.Println("Initializing submodules using Git CLI...")
	if output, err := RunGitCommand("submodule", "init"); err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, output)
	} else {
		fmt.Println("Submodules initialized successfully using Git CLI.")
	}

	// Update submodules
	fmt.Println("Updating submodules using Git CLI...")
	if output, err := RunGitCommand("submodule", "update"); err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, output)
	} else {
		fmt.Println("Submodules updated successfully using Git CLI.")
	}
}

func notMain() {
	// Replace these variables with your own values
	localPath := "/home/rnm/Dev/versionctrls/gitcliexperiments/repotest"
	// repoURL := "https://github.com/renatonmag/gitcliexperiments.git"
	// submoduleURL := "https://github.com/renatonmag/gitexperimentsintegration.git"
	// submodulePath := "submodule-directory"
	// username := "renatonmag"
	// password := "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"

	// // Clone the repository using go-git
	// cloneRepoWithGoGit(repoURL, localPath, username, password)

	// // Add a submodule using the Git CLI
	// addSubmoduleWithGitCLI(localPath, submodulePath, submoduleURL)

	// Replace these variables with your own values
	commitMessage := "Add line 2"
	authorName := "Renato Nankran"
	authorEmail := "renato.n@example.com"

	// Commit changes
	if err := commitChanges(localPath, commitMessage, authorName, authorEmail); err != nil {
		log.Fatalf("Failed to commit changes: %v", err)
	}
}

// package main

// import (
//     "fmt"
//     "log"

//     "gopkg.in/src-d/go-git.v4"
//     "gopkg.in/src-d/go-git.v4/plumbing"
//     "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
// )

// func main() {
//     r, err := git.PlainOpen("<REPOSITORY_PATH>")
//     if err != nil {
//         log.Fatal(err)
//     }

//     branch := fmt.Sprintf("refs/heads/%s", "master")

//     w, err := r.Worktree()
//     if err != nil {
//         log.Fatal(err)
//     }

//     if err := w.Pull(&git.PullOptions{
//         ReferenceName: plumbing.ReferenceName(branch),
//         Auth: &http.BasicAuth{
//             Username: "<GITHUB_USERNAME>",
//             Password: "<GITHUB_API_KEY>",
//         },
//     }); err != nil {
//         log.Fatal(err)
//     }
// }
