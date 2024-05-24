package repository

import (
	"fmt"
	"os"
	"os/exec"
)

// RunGitCommand runs a Git command and returns the output
func (r Repository) RunGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run git %v: %v\nOutput: %s", args, err, string(output))
	}
	return string(output), nil
}
