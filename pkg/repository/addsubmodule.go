package repository

import (
	"fmt"
	"log"
)

func (r Repository) AddSubmodule(url, path string) error {
	fmt.Println("Adding submodule...")
	output, err := r.RunGitCommand("submodule", "add", url, path)
	if err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, output)
		return err
	} else {
		fmt.Println("Submodules added successfully.")
	}

	// Initialize submodules
	fmt.Println("Initializing submodules...")
	if output, err := r.RunGitCommand("submodule", "init"); err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, output)
		return err
	} else {
		fmt.Println("Submodules initialized successfully.")
	}

	// Update submodules
	fmt.Println("Updating submodules...")
	if output, err := r.RunGitCommand("submodule", "update"); err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, output)
		return err
	} else {
		fmt.Println("Submodules updated successfully.")
	}

	return nil
}
