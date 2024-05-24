package gitmanagement

import (
	"fmt"
	"log"
	"os"
)

// Add a submodule using the Git CLI
func AddSubmodule(repoPath, submodulePath, submoduleURL string) {
	// Change directory to the cloned repository
	if err := os.Chdir(repoPath); err != nil {
		log.Fatalf("Failed to change directory: %v", err)
	}

	// Add a submodule
	fmt.Println("Adding submodule...")
	if output, err := RunGitCommand("submodule", "add", submoduleURL, submodulePath); err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, output)
	} else {
		fmt.Println("Submodule added successfully.")
	}

	// Initialize submodules
	fmt.Println("Initializing submodules...")
	if output, err := RunGitCommand("submodule", "init"); err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, output)
	} else {
		fmt.Println("Submodules initialized successfully.")
	}

	// Update submodules
	fmt.Println("Updating submodules...")
	if output, err := RunGitCommand("submodule", "update"); err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, output)
	} else {
		fmt.Println("Submodules updated successfully.")
	}
}
