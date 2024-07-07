package installer

import "os/exec"

func InstallToGo() error {
	return nil
}

func DeployToServer() error {
	return nil
}

// isCommandAvailable checks if a command is available in the system
func IsCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
