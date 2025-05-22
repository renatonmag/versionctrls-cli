package watching

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func WatchDir(dir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	
	