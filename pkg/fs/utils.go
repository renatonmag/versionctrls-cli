package fs

import (
	"fmt"
	"os"
	"syscall"
)

func GetInode(path string) (uint64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("failed to cast file info to Stat_t")
	}

	return stat.Ino, nil
}
