package npm

import (
	"fmt"
	"runtime"

	"github.com/theHamdiz/gost/installer"
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
	// Check the operating system
	osType := runtime.GOOS
	fmt.Println(">>Gost>> Operating system detected:", osType)

	// Check if Node.js is installed
	if !installer.IsCommandAvailable("node") {
		fmt.Println(">>Gost>> Node.js is not installed. Installing Node.js...")
		if err := installNodeJS(osType); err != nil {
			return err
		}
	} else {
		fmt.Println(">>Gost>> Node.js is already installed.")
	}

	// Install npm
	fmt.Println(">>Gost>> Installing npm...")
	if osType == "windows" {
		if err := runner.RunCommand("cmd", "/C", "npm install -g npm"); err != nil {
			return err
		}
	} else {
		if err := runner.RunCommand("sh", "-c", "curl -L https://www.npmjs.com/install.sh | sh"); err != nil {
			return err
		}
	}

	return nil
}

// installNodeJS installs the latest version of Node.js based on the operating system
func installNodeJS(osType string) error {
	var err error
	switch osType {
	case "windows":
		err = runner.RunCommand("powershell", "-Command", "Invoke-WebRequest -Uri https://nodejs.org/dist/latest/node-v18.17.1-x64.msi -OutFile nodejs.msi; Start-Process msiexec.exe -ArgumentList '/i', 'nodejs.msi', '/quiet', '/norestart' -NoNewWindow -Wait")
	case "darwin":
		err = runner.RunCommand("sh", "-c", "curl -fsSL https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash && source ~/.nvm/nvm.sh && nvm install node")
	case "linux":
		err = runner.RunCommand("sh", "-c", "curl -fsSL https://deb.nodesource.com/setup_current.x | sudo -E bash - && sudo apt-get install -y nodejs")
	default:
		return fmt.Errorf("unsupported operating system: %s", osType)
	}
	return err
}
