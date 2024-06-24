package runner

import (
	"fmt"
	"os"
	"os/exec"
)

// IsBinaryInstalled checks if a given binary is installed and callable on the system
func IsBinaryInstalled(binary string) bool {
	path, err := exec.LookPath(binary)
	if err != nil {
		return false
	}
	fmt.Printf("%s is available at %s\n", binary, path)
	return true
}

func RunCommand(name string, arg ...string) error {
	return RunCommandWithDir("", name, arg...)
}

func RunCommandWithDir(dir string, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func RunProject(projectDir string) error {
	return RunCommandWithDir(projectDir, "go", "run", "cmd/app/main.go")
}
