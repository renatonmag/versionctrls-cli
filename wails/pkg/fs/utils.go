package fs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

func IsFileTooLarge(path string, size int64) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Size() >= size
}

func GetFile(repoPath string, fileName string) string {
	return filepath.Join(repoPath, fileName)
}

func AppendToGitignore(gitignorePath, relPath string) error {
	// First, check if relPath is already present in the file
	if data, err := os.ReadFile(gitignorePath); err == nil {
		lines := string(data)
		for _, line := range splitLines(lines) {
			if line == relPath {
				// Already present, nothing to do
				return nil
			}
		}
	}
	// Not present, append it
	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(relPath + "\n")
	return err
}

// splitLines splits a string into lines, handling both \n and \r\n
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			// Remove trailing \r if present (for Windows line endings)
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	// Add last line if not empty
	if start < len(s) {
		line := s[start:]
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		lines = append(lines, line)
	}
	return lines
}

// LineExistsInFile returns true if any line in filePath is exactly equal to target.
func LineExistsInFile(filePath, target string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == target {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}

// RemoveLineFromFile removes all lines exactly equal to the given string from the file at filePath.
func RemoveLineFromFile(filePath, target string) error {
	// Read all lines
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != target {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Write filtered lines back to the file
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, line := range lines {
		_, err := outFile.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
