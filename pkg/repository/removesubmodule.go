package repository

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// RemoveSubmodule removes a submodule from the repository
// func (r *Repository) RemoveSubmodule() error {
// 	if r.repo == nil {
// 		return errors.New("no repository opened")
// 	}

// 	// Open the worktree
// 	worktree, err := r.repo.Worktree()
// 	if err != nil {
// 		return err
// 	}

// 	// Remove the submodule from .gitmodules
// 	repoRoot, err := r.GetRepoRoot()
// 	if err != nil {
// 		return err
// 	}
// 	gitmodulesPath := filepath.Join(repoRoot, ".gitmodules")
// 	cfg, err := ini.Load(gitmodulesPath)
// 	if err != nil {
// 		log.Fatal("Failed to load .gitmodules")
// 		return err
// 	}

// 	section := cfg.Section(fmt.Sprintf("submodule \"%s\"", r.submodulePath))
// 	if section != nil {
// 		cfg.DeleteSection(section.Name())
// 		err = cfg.SaveTo(gitmodulesPath)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// Remove the submodule from .git/config
// 	configPath := filepath.Join(repoRoot, ".git", "config")
// 	configCfg, err := ini.Load(configPath)
// 	if err != nil {
// 		log.Fatal("Failed to load .git/config")
// 		return err
// 	}

// 	section = configCfg.Section(fmt.Sprintf("submodule \"%s\"", r.submodulePath))
// 	if section != nil {
// 		configCfg.DeleteSection(section.Name())
// 		err = configCfg.SaveTo(configPath)
// 		if err != nil {
// 			log.Fatal("Failed to save .git/config")
// 			return err
// 		}
// 	}

// 	// Remove the submodule directory from the Git index
// 	_, err = worktree.Remove(r.submodulePath)
// 	if err != nil {
// 		log.Fatal("Failed to remove submodule from index")
// 		return err
// 	}

// 	// Delete the submodule directory from the working directory
// 	submoduleDir := filepath.Join(repoRoot, r.submodulePath)
// 	err = os.RemoveAll(submoduleDir)
// 	if err != nil {
// 		log.Fatal("Failed to remove submodule directory")
// 		return err
// 	}

// 	return nil
// }

func (r *Repository) RemoveSubmodule() error {
	if r.repo == nil {
		return errors.New("no repository opened")
	}
	repoRoot, _ := r.GetRepoRoot()

	// Remove the submodule entry from the .git/config.
	configPath := filepath.Join(repoRoot, ".git", "config")
	cfg, err := ini.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load git config: %w", err)
	}

	submoduleSection := fmt.Sprintf("submodule.%s", r.submodulePath)
	cfg.DeleteSection(submoduleSection)

	if err := cfg.SaveTo(configPath); err != nil {
		return fmt.Errorf("failed to save git config: %w", err)
	}

	// Remove the submodule directory from .git/modules.
	gitModulesPath := filepath.Join(repoRoot, ".git", "modules", r.submodulePath)
	if err := os.RemoveAll(gitModulesPath); err != nil {
		return fmt.Errorf("failed to remove submodule directory from .git/modules: %w", err)
	}

	gitmodulesPath := filepath.Join(repoRoot, ".gitmodules")
	cfg, err = ini.Load(gitmodulesPath)
	if err != nil {
		log.Fatal("Failed to load .gitmodules")
		return err
	}

	section := cfg.Section(fmt.Sprintf("submodule \"%s\"", r.submodulePath))
	if section != nil {
		cfg.DeleteSection(section.Name())
		err = cfg.SaveTo(gitmodulesPath)
		if err != nil {
			return err
		}
	}

	// Remove the submodule directory from the repository.
	submoduleFullPath := filepath.Join(repoRoot, r.submodulePath)
	if err := os.RemoveAll(submoduleFullPath); err != nil {
		return fmt.Errorf("failed to remove submodule directory from repository: %w", err)
	}

	return nil
}
