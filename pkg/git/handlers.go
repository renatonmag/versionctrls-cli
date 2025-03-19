package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/helshabini/fsbroker"
	"github.com/renatonmag/version-ctrls-cli/pkg/fs"
)

func OnCreate(repo *Repository, event *fsbroker.FSEvent) error {
	branchName := fmt.Sprintf("(%s)-(%s)", GetShortName(repo.MainRepoBranch), event.Path)

	exists := repo.BranchExists(branchName)
	if !exists {
		err := repo.CreateBranch(branchName, "master")
		if err != nil {
			return err
		}
	}

	err := repo.CheckoutBranch(branchName)
	if err != nil {
		return err
	}

	err = fs.NewFsService("").Replicate.CreateHardlink(event.Path, repo.Path)
	if err != nil {
		return err
	}

	// commit change
	file := filepath.Base(event.Path)
	hash, err := repo.CommitToBranch(branchName, file, fmt.Sprintf("File created: %s", event.Path))
	if err != nil {
		return err
	}
	fmt.Println("hash: ", hash)
	return nil
}

func OnRemove(repo *Repository, event *fsbroker.FSEvent) error {
	file := filepath.Base(event.Path)
	branchName := fmt.Sprintf("(%s)-(%s)", GetShortName(repo.MainRepoBranch), event.Path)

	exists := repo.BranchExists(branchName)
	if !exists {
		return fmt.Errorf("branch %s does not exist", branchName)
	}

	err := repo.CheckoutBranch(branchName)
	if err != nil {
		return err
	}

	destPath := filepath.Join(repo.Path, file)
	err = os.Remove(destPath)
	if err != nil {
		fmt.Println("error removing file: ", err)
		return err
	}
	fmt.Println("destPath: ", destPath)

	// commit change
	hash, err := repo.CommitToBranch(branchName, file, fmt.Sprintf("File deleted: %s", event.Path))
	if err != nil {
		return err
	}
	fmt.Println("hash: ", hash)

	return nil
}

// func onModify(branchName, srcBranch string) error {

// }

// func onRename(branchName, srcBranch string) error {

// }

// func onMove(branchName, srcBranch string) error {

// }
