package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/renatonmag/fsbroker"
	"github.com/renatonmag/version-ctrls/pkg/fs"
	"github.com/renatonmag/version-ctrls/pkg/utils"
)

func OnCreate(repo *LocalRepository, event *fsbroker.FSEvent) error {
	branchName := getBranchName(event.Path)
	exists := repo.BranchExists(branchName)
	max_file_size := repo.AppConfig.Integration.MaxFileSize
	integrationPath := repo.AppConfig.Integration.Path
	gitignorePath := fs.GetFile(integrationPath, ".gitignore")
	if fs.IsFileTooLarge(event.Path, max_file_size) && !exists {
		fmt.Printf("file %s is too large, skipping\n", event.Path)
		fmt.Printf("Adding to ignore list\n")
		fs.AppendToGitignore(gitignorePath, event.Path)
		return nil
	} else if exists, _ := fs.LineExistsInFile(gitignorePath, event.Path); exists {
		fs.RemoveLineFromFile(gitignorePath, event.Path)
	}

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
	baseFilePath := filepath.Join("base", event.Path)
	_fileToCreateAtRepo := filepath.Join(repo.Path, baseFilePath)
	_dirPath := filepath.Dir(_fileToCreateAtRepo)
	exists = utils.DirectoryExists(_fileToCreateAtRepo)
	if !exists {
		fmt.Println("directory does not exist, creating it")
		err = os.MkdirAll(_dirPath, 0755)
		if err != nil {
			return err
		}
	}

	err = fs.NewFsService().Replicate.CreateHardlink(event.Path, _dirPath)
	if err != nil {
		fmt.Println("error creating hardlink: ", err)
		return err
	}

	// commit change
	_, err = repo.CommitToBranch(branchName, baseFilePath, fmt.Sprintf("File created: %s", event.Path), false)
	if err != nil {
		fmt.Println("error committing to branch: ", err)
		return err
	}
	return nil
}

func OnRemove(repo *LocalRepository, event *fsbroker.FSEvent) error {
	file := filepath.Base(event.Path)
	branchName := getBranchName(event.Path)

	exists := repo.BranchExists(branchName)
	if !exists {
		return fmt.Errorf("branch %s does not exist", branchName)
	}

	err := repo.CheckoutBranch(branchName)
	if err != nil {
		return err
	}

	// commit change
	hash, err := repo.CommitToBranch(branchName, file, fmt.Sprintf("File deleted: %s", event.Path), false)
	if err != nil {
		return err
	}
	fmt.Println("hash: ", hash)

	return nil
}

func OnModify(repo *LocalRepository, event *fsbroker.FSEvent) error {
	branchName := getBranchName(event.Path)
	exists := repo.BranchExists(branchName)
	if !exists {
		err := OnCreate(repo, event)
		if err != nil {
			fmt.Println("error creating branch: ", err)
			return err
		}
		return nil
	}

	baseFilePath := filepath.Join("base", event.Path)
	_fileToCreateAtRepo := filepath.Join(repo.Path, baseFilePath)
	_dirPath := filepath.Dir(_fileToCreateAtRepo)

	originalInode, err := fs.GetInode(event.Path)
	if err != nil {
		return err
	}
	remoteInode, err := fs.GetInode(_fileToCreateAtRepo)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file does not exist, treat as not an error for inode comparison
			remoteInode = 0
		} else {
			return err
		}
	}
	isSameFile := originalInode == remoteInode

	if !isSameFile {
		err = repo.CheckoutBranch(branchName)
		if err != nil {
			return err
		}

		exists = utils.DirectoryExists(_fileToCreateAtRepo)
		if !exists {
			fmt.Println("directory does not exist, creating it")
			err = os.MkdirAll(_dirPath, 0755)
			if err != nil {
				return err
			}
		}

		if utils.FileExists(_fileToCreateAtRepo) {
			err := os.Remove(_fileToCreateAtRepo)
			if err != nil {
				return err
			}
		}
		err = fs.NewFsService().Replicate.CreateHardlink(event.Path, _dirPath)
		if err != nil {
			fmt.Println("error creating hardlink: ", err)
			return err
		}
	}

	// if fileExists(event.Path) {
	// 	os.Remove(event.Path)
	// 	err = fs.NewFsService("").Replicate.CreateHardlink(event.Path, _dirPath)
	// 	if err != nil {
	// 		fmt.Println("error creating hardlink: ", err)
	// 		return err
	// 	}
	// }

	// commit change
	_, err = repo.CommitToBranch(branchName, baseFilePath, fmt.Sprintf("File modified: %s", event.Path), false)
	if err != nil {
		fmt.Println("error committing to branch: ", err)
		return err
	}
	return nil

}

// func onRename(branchName, srcBranch string) error {

// }

func OnMove(repo *LocalRepository, event *fsbroker.FSEvent) error {
	oldPath := event.Properties["OldPath"]
	branchName := getBranchName(oldPath)
	newBranchName := getBranchName(event.Path)
	exists := repo.BranchExists(branchName)
	if exists {
		err := repo.RenameBranch(branchName, newBranchName)
		if err != nil {
			fmt.Println("error renaming branch: ", err)
			return err
		}
		// err = repo.CheckoutBranch(newBranchName)
		// if err != nil {
		// 	fmt.Println("error checking out branch: ", err)
		// 	return err
		// }
	} else {
		err := OnCreate(repo, event)
		if err != nil {
			fmt.Println("error creating branch: ", err)
			return err
		}
		return nil
	}

	err := repo.CheckoutBranch(newBranchName)
	if err != nil {
		return err
	}

	baseFilePath := filepath.Join("base", event.Path)
	_fileToCreateAtRepo := filepath.Join(repo.Path, baseFilePath)
	_dirPath := filepath.Dir(_fileToCreateAtRepo)

	exists = utils.DirectoryExists(_fileToCreateAtRepo)
	if !exists {
		fmt.Println("directory does not exist, creating it")
		err = os.MkdirAll(_dirPath, 0755)
		if err != nil {
			return err
		}
	}

	if utils.FileExists(_fileToCreateAtRepo) {
		err := os.Remove(_fileToCreateAtRepo)
		if err != nil {
			return err
		}
	}
	err = fs.NewFsService().Replicate.CreateHardlink(event.Path, _dirPath)
	if err != nil {
		fmt.Println("error creating hardlink: ", err)
		return err
	}

	// tryCount := 0
	// maxTries := 3
	// for {
	// 	if ok := fileExists(event.Path); ok {
	// 		fmt.Printf("file exists")
	// 		break // file exists
	// 	} else {
	// 		tryCount++
	// 		if tryCount >= maxTries {
	// 			return fmt.Errorf("file %s does not exist after %d tries", event.Path, maxTries)
	// 		}
	// 		time.Sleep(50 * time.Millisecond)
	// 	}
	// }

	_, err = repo.CommitToBranch(newBranchName, baseFilePath, fmt.Sprintf("File moved: %s", event.Path), true)
	if err != nil {
		fmt.Println("error committing to branch: ", err)
		return err
	}
	return nil
}
