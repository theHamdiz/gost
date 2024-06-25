package npm

import (
	"fmt"

	"github.com/theHamdiz/gost/runner"
)

// CheckNPMInstalled checks whether npm is installed, if not attempts to install it and checks again
func CheckNPMInstalled() error {
	// Check if npm is installed
	if err := isNPMInstalled(); err == nil {
		return nil
	}

	fmt.Println(">>Gost>> npm not found, attempting to install...")

	// Attempt to install npm
	if err := installNPM(); err != nil {
		return fmt.Errorf(">>Gost>> failed to install npm: %w", err)
	}

	// Check again if npm is installed
	if err := isNPMInstalled(); err != nil {
		return fmt.Errorf(">>Gost>> npm installation verification failed: %w", err)
	}

	return nil
}

// isNPMInstalled checks if npm is installed
func isNPMInstalled() error {
	if err := runner.RunCommand("npm", "--version"); err != nil {
		return err
	}
	return nil
}

// installNPM attempts to install npm
func installNPM() error {
	// Using the appropriate command for the operating system
	if err := runner.RunCommand("sh", "-c", "curl -L https://www.npmjs.com/install.sh | sh"); err != nil {
		return err
	}
	return nil
}
