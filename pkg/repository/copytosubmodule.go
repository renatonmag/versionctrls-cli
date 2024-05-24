package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/renatonmag/versionctrls-cli/pkg/utils"
)

const maxFileSize = 500 * 1024 // 500KB in bytes

// CopyChangedFilesToSubmodule copies all changed files from the root repository to the submodule
func (r Repository) CopyChangedFilesToSubmodule() error {
	changedFiles, err := r.GetChangedFiles()
	if err != nil {
		return err
	}

	rootPath, _ := r.GetRepoRoot()

	for _, file := range changedFiles {
		srcPath := filepath.Join(rootPath, file)
		dstPath := filepath.Join(r.submodulePath, file)

		fileInfo, err := os.Stat(srcPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Skipping %s as it does not exist\n", srcPath)
				continue
			}
			return err
		}

		if fileInfo.Size() > maxFileSize {
			fmt.Printf("Skipping %s as it exceeds the file size limit of 500KB\n", srcPath)
			continue
		}

		// Ensure the destination directory exists
		dstDir := filepath.Dir(dstPath)
		err = os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return err
		}

		err = utils.CopyFile(srcPath, dstPath)
		if err != nil {
			return err
		}
		fmt.Printf("Copied %s to %s\n", srcPath, dstPath)
	}

	return nil
}
