package utils

import (
	"fmt"
	"os"
)

func DirectoryExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		// File or directory exists
		fileInfo, err := os.Stat(path)
		if err != nil {
			// Handle the error of getting file info (unlikely, but possible)
			fmt.Println("Error getting file info:", err) // Log or handle appropriately
			return false                                 // Indicate failure, don't assume existence
		}
		return fileInfo.IsDir() // Check if it's a directory
	} else if os.IsNotExist(err) {
		// Directory does not exist
		return false
	} else {
		// Some other error occurred (e.g., permission denied)
		fmt.Println("Error checking directory:", err) // Log or handle appropriately
		return false                                  // Indicate failure, don't assume existence
	}
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true // file exists
	} else if os.IsNotExist(err) {
		return false
	} else {
		fmt.Printf("error checking file existence: %w", err)
	}
	return false
}
