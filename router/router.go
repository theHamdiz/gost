package router

import (
	"os"
	"path/filepath"
)

func GetCWD() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return cwd, nil
}

func GetProjectPath(appName string) (string, error) {
	cwd, err := GetCWD()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, appName), nil
}

func GetDbPath(appName string) (string, error) {
	project, err := GetProjectPath(appName)
	if err != nil {
		return "", err
	}
	return filepath.Join(project, "app", "db", "data.db"), nil
}
