package fs

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/renatonmag/version-ctrls/pkg/config"
	"github.com/renatonmag/version-ctrls/pkg/utils"
	ignore "github.com/sabhiram/go-gitignore"
)

type FsService struct {
	Replicate *ReplicateDir
}

type IgnoreService struct {
	filePath  string
	matcher   *ignore.GitIgnore
	appConfig *config.ApplicationConfig
}

func NewFsService() *FsService {
	return &FsService{
		Replicate: &ReplicateDir{
			ignore: NewIgnoreService(),
		},
	}
}

func NewIgnoreService() *IgnoreService {
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
	}
	ignore := &IgnoreService{
		appConfig: appConfig,
	}
	ignoreExists := utils.FileExists(ignore.filePath)
	if !ignoreExists {
		ignore.compileIgnoreLines(ignore.filePath)
	}
	return ignore
}

func (i *IgnoreService) compileIgnoreLines(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		// If the file doesn't exist, use an empty matcher
		i.matcher = ignore.CompileIgnoreLines()
		return
	}
	defer file.Close()

	var lines []string
	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	i.matcher = ignore.CompileIgnoreLines(lines...)
}

func (i *IgnoreService) MatchesPath(path string) bool {
	return i.matcher.MatchesPath(path)
}

type ReplicateDir struct {
	ignore *IgnoreService
}

// func (r *ReplicateDir) SyncDirs(src, dst string) error {
// 	diffs, err := r.DiffDirs(src, dst)
// 	if err != nil {
// 		return err
// 	}

// 	err = r.CreateHardlinks(diffs[src], dst)
// 	if err != nil {
// 		return err
// 	}

// 	err = r.CleanFiles(diffs[dst], dst)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *ReplicateDir) CleanFiles(paths []string, dst string) error {
	for _, path := range paths {
		err := os.Remove(filepath.Join(dst, path))
		if err != nil {
			return err
		}
	}

	return nil
}

// func (r *ReplicateDir) CreateHardlinks(paths []string, dst string) error {
// 	for _, path := range paths {
// 		// Get file info
// 		info, err := os.Stat(path)
// 		if err != nil {
// 			return err
// 		}

// 		destPath := filepath.Join(dst, path)

// 		if info.IsDir() {
// 			// Create the directory structure
// 			if err := os.MkdirAll(destPath, 0755); err != nil {
// 				return err
// 			}
// 		} else {
// 			// Create parent directories if they don't exist
// 			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
// 				return err
// 			}
// 			// For files, create a hard link
// 			if err := os.Link(path, destPath); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

func (r *ReplicateDir) CreateHardlink(path, dst string) error {
	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Get the last part of the path (filename)
	filename := filepath.Base(path)
	destPath := filepath.Join(dst, filename)

	if info.IsDir() {
		return fmt.Errorf("cannot create hardlink for directory: %s", path)
	}

	// For files, create a hard link
	if err := os.Link(path, destPath); err != nil {
		return err
	}
	chmod := os.Chmod(destPath, info.Mode())
	if chmod != nil {
		return chmod
	}
	return nil
}

// CleanWorkingTree removes all contents of a directory except for items in the ignore list.
// The directory itself is preserved. Items in the ignore slice are directory or file names
// (not paths) that should be preserved.
func (r *ReplicateDir) CleanWorkingTree(dir string) error {

	dotgit := ".git"

	// Read directory contents
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Process each entry
	for _, entry := range entries {
		name := entry.Name()

		// Skip if in ignore list
		if name == dotgit {
			continue
		}

		// Build full path
		path := filepath.Join(dir, name)

		// Remove the item (recursively if it's a directory)
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}
	}

	return nil
}

type FileMeta struct {
	Size int64
	Mode fs.FileMode
}

// buildFileMap recursively builds a map of relative file paths to their metadata.
// Files matching patterns in gitignore are excluded.
func (r *ReplicateDir) BuildFileMap(root string) (map[string]FileMeta, error) {
	m := make(map[string]FileMeta)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		// Skip the root itself.
		if rel == "." {
			return nil
		}

		// Skip files that match gitignore patterns
		if r.ignore != nil && r.ignore.MatchesPath(rel) {
			if d.IsDir() {
				return filepath.SkipDir // Skip this directory and its contents
			}
			return nil // Skip this file
		}

		if d.IsDir() {
			// Optionally record directories, if desired.
			m[rel] = FileMeta{}
		} else {
			info, err := d.Info()
			if err != nil {
				return err
			}

			m[rel] = FileMeta{
				Size: info.Size(),
				Mode: info.Mode(),
			}
		}
		return nil
	})
	return m, err
}

// DiffDirs compares two directory trees and returns differences organized by directory.
// Returns a map with directory paths as keys and slices of unique files as values.
func (r *ReplicateDir) DiffDirs(dir1, dir2 string) (map[string][]string, error) {
	map1, err := r.BuildFileMap(dir1)
	if err != nil {
		return nil, err
	}
	map2, err := r.BuildFileMap(dir2)
	if err != nil {
		return nil, err
	}

	// Create result structure
	result := map[string][]string{
		dir1: {},
		dir2: {},
	}

	// Check for files only in dir1 or differences
	for rel := range map1 {
		if _, ok := map2[rel]; !ok {
			// File only exists in dir1
			result[dir1] = append(result[dir1], rel)
		}
	}

	// Check for files only in dir2
	for rel := range map2 {
		if _, ok := map1[rel]; !ok {
			// File only exists in dir2
			result[dir2] = append(result[dir2], rel)
		}
	}

	return result, nil
}

type FsWatcher struct {
	watcher *fsnotify.Watcher
}

func NewFsWatcher() (*FsWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &FsWatcher{
		watcher: watcher,
	}, nil
}

func (r *FsWatcher) WatchDir(dir string) error {
	// Add a path.
	err := r.watcher.Add(dir)
	if err != nil {
		return err
	}

	// Start listening for events.
	go func() {
		// Create a map to track recent write events by filename
		writeEvents := make(map[string]bool)

		// Create a timer for coalescing events
		var timer *time.Timer

		// Function to process accumulated events
		processEvents := func() {
			// Process all accumulated write events
			for filename := range writeEvents {
				fmt.Println("modified file:", filename)
			}
			// Clear the map after processing
			writeEvents = make(map[string]bool)
		}

		for {
			select {
			case event, ok := <-r.watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)

				if event.Has(fsnotify.Write) {
					// Instead of printing immediately, add to our tracking map
					writeEvents[event.Name] = true

					// Reset or create the coalescing timer
					if timer != nil {
						timer.Stop()
					}
					timer = time.AfterFunc(300*time.Millisecond, processEvents)
				}

				// Handle other event types immediately
				if event.Has(fsnotify.Rename) {
					fmt.Println("renamed file Op:", event.Op)
					fmt.Println("renamed file:", event.Name)
				}
				if event.Has(fsnotify.Create) {
					fmt.Println("created file:", event.Name)
				}
				if event.Has(fsnotify.Remove) {
					fmt.Println("removed file:", event.Name)
				}
			case err, ok := <-r.watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	return nil
}
