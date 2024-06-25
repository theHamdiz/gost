package cfg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
)

func createTestFile(t *testing.T, filePath string, content string) {
	file, err := os.Create(filePath)
	assert.NoError(t, err)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(">>Gost>>Error closing file:", err)
		}
	}(file)

	_, err = file.WriteString(content)
	assert.NoError(t, err)
}

func TestSaveAsEnv(t *testing.T) {
	config := &GostConfig{
		PreferredIDE:                 "vscode",
		PreferredBackendFramework:    "echo",
		PreferredUiFramework:         "react",
		PreferredComponentsFramework: "component-framework",
		PreferredDbDriver:            "postgres",
		PreferredDbOrm:               "gorm",
		PreferredPort:                8080,
		GlobalSettings:               "global-settings",
		PreferredConfigFormat:        "env",
	}

	filePath := "test.env"
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println(">>Gost>> Error removing file:", err)
		}
	}(filePath)

	err := config.SaveAsEnv(filePath)
	assert.NoError(t, err)

	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	expected := "PreferredIDE=vscode\nPreferredBackendFramework=echo\nPreferredUiFramework=react\nPreferredComponentsFramework=component-framework\nPreferredDbDriver=postgres\nPreferredDbOrm=gorm\nPreferredPort=8080\nGlobalSettings=global-settings\nPreferredConfigFormat=env\n"
	assert.Equal(t, expected, string(content))
}

func TestSaveAsJSON(t *testing.T) {
	config := &GostConfig{
		PreferredIDE:                 "vscode",
		PreferredBackendFramework:    "echo",
		PreferredUiFramework:         "react",
		PreferredComponentsFramework: "component-framework",
		PreferredDbDriver:            "postgres",
		PreferredDbOrm:               "gorm",
		PreferredPort:                8080,
		GlobalSettings:               "global-settings",
		PreferredConfigFormat:        "json",
	}

	filePath := "test.json"
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println(">>Gost>>Error removing file:", err)
		}
	}(filePath)

	err := config.SaveAsJSON(filePath)
	assert.NoError(t, err)

	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	expected, err := json.MarshalIndent(config, "", "  ")
	assert.NoError(t, err)
	assert.JSONEq(t, string(expected), string(content))
}

func TestSaveAsTOML(t *testing.T) {
	config := &GostConfig{
		PreferredIDE:                 "vscode",
		PreferredBackendFramework:    "echo",
		PreferredUiFramework:         "react",
		PreferredComponentsFramework: "component-framework",
		PreferredDbDriver:            "postgres",
		PreferredDbOrm:               "gorm",
		PreferredPort:                8080,
		GlobalSettings:               "global-settings",
		PreferredConfigFormat:        "toml",
	}

	filePath := "test.toml"
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println(">>Gost>>Error removing file:", err)
		}
	}(filePath)

	err := config.SaveAsTOML(filePath)
	assert.NoError(t, err)

	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	expected, err := toml.Marshal(config)
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(content))
}

func TestLoadFromEnv(t *testing.T) {
	content := `PreferredIDE=vscode
PreferredBackendFramework=echo
PreferredUiFramework=react
PreferredComponentsFramework=component-framework
PreferredDbDriver=postgres
PreferredDbOrm=gorm
PreferredPort=8080
GlobalSettings=global-settings
PreferredConfigFormat=env
`
	filePath := "test.env"
	createTestFile(t, filePath, content)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println(">>Gost>>Error removing file:", err)
		}
	}(filePath)

	config, err := LoadFromEnv(filePath)
	assert.NoError(t, err)

	expected := &GostConfig{
		PreferredIDE:                 "vscode",
		PreferredBackendFramework:    "echo",
		PreferredUiFramework:         "react",
		PreferredComponentsFramework: "component-framework",
		PreferredDbDriver:            "postgres",
		PreferredDbOrm:               "gorm",
		PreferredPort:                8080,
		GlobalSettings:               "global-settings",
		PreferredConfigFormat:        "env",
	}

	assert.Equal(t, expected, config)
}

func TestLoadFromJSON(t *testing.T) {
	config := &GostConfig{
		PreferredIDE:                 "vscode",
		PreferredBackendFramework:    "echo",
		PreferredUiFramework:         "react",
		PreferredComponentsFramework: "component-framework",
		PreferredDbDriver:            "postgres",
		PreferredDbOrm:               "gorm",
		PreferredPort:                8080,
		GlobalSettings:               "global-settings",
		PreferredConfigFormat:        "json",
	}
	filePath := "test.json"
	content, err := json.MarshalIndent(config, "", "  ")
	assert.NoError(t, err)
	createTestFile(t, filePath, string(content))
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println(">>Gost>>Error removing file:", err)
		}
	}(filePath)

	loadedConfig, err := LoadFromJSON(filePath)
	assert.NoError(t, err)
	assert.Equal(t, config, loadedConfig)
}

func TestLoadFromTOML(t *testing.T) {
	config := &GostConfig{
		PreferredIDE:                 "vscode",
		PreferredBackendFramework:    "echo",
		PreferredUiFramework:         "react",
		PreferredComponentsFramework: "component-framework",
		PreferredDbDriver:            "postgres",
		PreferredDbOrm:               "gorm",
		PreferredPort:                8080,
		GlobalSettings:               "global-settings",
		PreferredConfigFormat:        "toml",
	}
	filePath := "test.toml"
	content, err := toml.Marshal(config)
	assert.NoError(t, err)
	createTestFile(t, filePath, string(content))
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println(">>Gost>>Error removing file:", err)
		}
	}(filePath)

	loadedConfig, err := LoadFromTOML(filePath)
	assert.NoError(t, err)
	assert.Equal(t, config, loadedConfig)
}

func TestGetIDEBinaryName(t *testing.T) {
	config := &GostConfig{
		PreferredIDE: "vscode",
	}

	binaryName, err := config.GetIDEBinaryName()
	assert.NoError(t, err)
	assert.Equal(t, "code", binaryName)

	config.PreferredIDE = "unknownIDE"
	_, err = config.GetIDEBinaryName()
	assert.Error(t, err)
	assert.Equal(t, "> Unknown IDE/editor", err.Error())
}

func TestResetConfig(t *testing.T) {
	content := `PreferredIDE=vscode
PreferredBackendFramework=echo
PreferredUiFramework=react
PreferredComponentsFramework=component-framework
PreferredDbDriver=postgres
PreferredDbOrm=gorm
PreferredPort=8080
GlobalSettings=global-settings
PreferredConfigFormat=env
`
	filePath := filepath.Join(os.TempDir(), ".gost.env")
	createTestFile(t, filePath, content)

	ResetConfig(filePath)

	_, err := os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))
}

func TestShowConfig(t *testing.T) {
	content := `PreferredIDE=vscode
PreferredBackendFramework=echo
PreferredUiFramework=react
PreferredComponentsFramework=component-framework
PreferredDbDriver=postgres
PreferredDbOrm=gorm
PreferredPort=8080
GlobalSettings=global-settings
PreferredConfigFormat=env
`
	filePath := filepath.Join(os.TempDir(), ".gost.env")
	createTestFile(t, filePath, content)

	ShowConfig(filePath)

	// Output is printed to the console; you can manually verify or capture output for automated testing
}
