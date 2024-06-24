package git

import (
	"fmt"

	"github.com/theHamdiz/gost/runner"
)

// CheckGitInstalled checks whether git is installed, if not attempts to install it and checks again
func CheckGitInstalled() error {
	// Check if git is installed
	if err := isGitInstalled(); err == nil {
		return nil
	}

	fmt.Println(">> git not found, attempting to install...")

	// Attempt to install git
	if err := installGit(); err != nil {
		return fmt.Errorf("failed to install git: %w", err)
	}

	// Check again if git is installed
	if err := isGitInstalled(); err != nil {
		return fmt.Errorf("git installation verification failed: %w", err)
	}

	return nil
}

// isGitInstalled checks if git is installed
func isGitInstalled() error {
	if err := runner.RunCommand("git", "--version"); err != nil {
		return err
	}
	return nil
}

// installGit attempts to install git
func installGit() error {
	// Using the appropriate command for the operating system
	if err := runner.RunCommand("sh", "-c", "sudo apt-get update && sudo apt-get install -y git"); err != nil {
		return err
	}
	return nil
}
