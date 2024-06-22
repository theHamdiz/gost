package gen

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml"
)

// Todo: refactor to GostConfig
type TemplateData struct {
	AppName                string
	BackendImport          string
	VersionedBackendImport string
	BackendInit            string
	ComponentsFramework    string
	Fingerprint            string
	UiFramework            string
	BackendPkg             string
	DbDriver               string
	DbOrm                  string
}

type GostConfig struct {
	PreferredBackendFramework    string
	PreferredUiFramework         string
	PreferredComponentsFramework string
	PreferredDbDriver            string
	PreferredDbOrm               string
	PreferredIDE                 string
	PreferredPort                int
	GlobalSettings               string
	PreferredConfigFormat        string
	AppName                      string
}

func (config *GostConfig) SaveAsEnv(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	fmt.Fprintln(writer, "PreferredIDE="+config.PreferredIDE)
	fmt.Fprintln(writer, "PreferredBackendFramework="+config.PreferredBackendFramework)
	fmt.Fprintln(writer, "PreferredUiFramework="+config.PreferredUiFramework)
	fmt.Fprintln(writer, "PreferredComponentsFramework="+config.PreferredComponentsFramework)
	fmt.Fprintln(writer, "PreferredDbDriver="+config.PreferredDbDriver)
	fmt.Fprintln(writer, "PreferredDbOrm="+config.PreferredDbOrm)
	fmt.Fprintf(writer, "PreferredPort=%d\n", config.PreferredPort)
	fmt.Fprintln(writer, "GlobalSettings="+config.GlobalSettings)
	fmt.Fprintln(writer, "PreferredConfigFormat="+config.PreferredConfigFormat)
	writer.Flush()

	return nil
}

func (config *GostConfig) SaveAsJSON(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

func (config *GostConfig) SaveAsTOML(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(config)
}

func LoadFromEnv(filePath string) (*GostConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &GostConfig{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		switch key {
		case "PreferredIDE":
			config.PreferredIDE = value
		case "PreferredBackendFramework":
			config.PreferredBackendFramework = value
		case "PreferredUiFramework":
			config.PreferredUiFramework = value
		case "PreferredComponentsFramework":
			config.PreferredComponentsFramework = value
		case "PreferredDbDriver":
			config.PreferredDbDriver = value
		case "PreferredDbOrm":
			config.PreferredDbOrm = value
		case "PreferredPort":
			config.PreferredPort, _ = strconv.Atoi(value)
		case "GlobalSettings":
			config.GlobalSettings = value
		case "PreferredConfigFormat":
			config.PreferredConfigFormat = value
		}
	}
	return config, scanner.Err()
}

func LoadFromJSON(filePath string) (*GostConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &GostConfig{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	return config, err
}

func LoadFromTOML(filePath string) (*GostConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &GostConfig{}
	decoder := toml.NewDecoder(file)
	err = decoder.Decode(config)
	return config, err
}

// GetBinaryName returns the binary file name for a given IDE/editor name.
func (c GostConfig) GetIDEBinaryName() (string, error) {
	ide := strings.ToLower(c.PreferredIDE)
	switch ide {
	case "vscode":
		return "code", nil
	case "goland":
		return "goland", nil
	case "idea":
		return "idea", nil
	case "cursor":
		return "cursor", nil
	case "zed":
		return "zed", nil
	case "sublime":
		return "subl", nil
	case "vim":
		return "vim", nil
	case "nvim":
		return "nvim", nil
	case "nano":
		return "nano", nil
	case "notepad++":
		return "notepad++", nil
	case "zeus":
		return "zeus", nil
	case "liteide":
		return "liteide", nil
	case "emacs":
		return "emacs", nil
	case "eclipse":
		return "eclipse", nil
	default:
		return "", errors.New("> Uknown IDE/editor")
	}
}
