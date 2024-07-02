package cfg

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/theHamdiz/gost/clr"
	"gopkg.in/yaml.v3"
)

type GostConfig struct {
	AppId                        string
	AppKey                       string
	AppName                      string
	DatabaseURI                  string
	EnvFile                      string
	EnvPassword                  string
	EnvUserName                  string
	GlobalSettings               string
	PreferredBackendFramework    string
	PreferredComponentsFramework string
	PreferredConfigFormat        string
	PreferredDbDriver            string
	PreferredDbOrm               string
	PreferredFrontEndFramework   string
	PreferredIDE                 string
	PreferredPort                int
	PreferredUiFramework         string
	RedisDb                      string
	RedisPassword                string
	RedisURI                     string
	Server                       string
}

type Configurable interface {
	SaveAsEnv(filePath string) error
	SaveAsJSON(filePath string) error
	SaveAsTOML(filePath string) error
	SaveAsYAML(filePath string) error
}

func (config *GostConfig) SaveAsEnv(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(">>Gost>> Error closing file:", err)
		}
	}(file)

	writer := bufio.NewWriter(file)
	_, err = fmt.Fprintln(writer, ">>Gost>> PreferredIDE="+config.PreferredIDE)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, ">>Gost>> PreferredBackendFramework="+config.PreferredBackendFramework)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, ">>Gost>> PreferredUiFramework="+config.PreferredUiFramework)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, ">>Gost>> PreferredComponentsFramework="+config.PreferredComponentsFramework)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, ">>Gost>> PreferredDbDriver="+config.PreferredDbDriver)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, ">>Gost>> PreferredDbOrm="+config.PreferredDbOrm)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(writer, ">>Gost>> PreferredPort=%d\n", config.PreferredPort)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, ">>Gost>> GlobalSettings="+config.GlobalSettings)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, ">>Gost>> PreferredConfigFormat="+config.PreferredConfigFormat)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, ">>Gost>> PreferredFrontEndFramework="+config.PreferredFrontEndFramework)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (config *GostConfig) SaveAsJSON(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(">>Gost>> Error closing file:", err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

func (config *GostConfig) SaveAsTOML(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(">>Gost>> Error closing file:", err)
		}
	}(file)

	encoder := toml.NewEncoder(file)
	return encoder.Encode(config)
}

func (config *GostConfig) SaveAsYAML(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(">>Gost>> Error closing file:", err)
		}
	}(file)

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(config)
}

// GetIDEBinaryName GetBinaryName returns the binary file name for a given IDE/editor name.
func (config *GostConfig) GetIDEBinaryName() (string, error) {
	ide := strings.ToLower(config.PreferredIDE)
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
		return "", errors.New(">>Gost>>  Unknown IDE/editor")
	}
}

func LoadFromEnv(filePath string) (*GostConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(">>Gost>> Error closing file:", err)
		}
	}(file)

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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(">>Gost>> Error closing file:", err)
		}
	}(file)

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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(">>Gost>> Error closing file:", err)
		}
	}(file)

	config := &GostConfig{}
	decoder := toml.NewDecoder(file)
	err = decoder.Decode(config)
	return config, err
}

func GetConfigPath() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(clr.Colorize(">>Gost>> Error getting user home directory", "red"))
		return ""
	}

	envFilePath := filepath.Join(usr.HomeDir, ".gost.env")
	jsonFilePath := filepath.Join(usr.HomeDir, ".gost.json")
	tomlFilePath := filepath.Join(usr.HomeDir, ".gost.toml")

	if _, err := os.Stat(envFilePath); err == nil {
		return envFilePath
	} else if _, err := os.Stat(jsonFilePath); err == nil {
		return jsonFilePath
	} else if _, err := os.Stat(tomlFilePath); err == nil {
		return tomlFilePath
	}
	return ""
}

func ResetConfig(configPath string) {
	if configPath == "" {
		configPath = GetConfigPath()
	}
	if configPath == "" {
		return
	}
	err := os.Remove(configPath)
	if err != nil {
		fmt.Println(clr.Colorize(">>Gost>> Error removing config file:", "red"))
		return
	}
}

func ShowConfig(configPath string) {
	if configPath == "" {
		configPath = GetConfigPath()
	}

	if configPath == "" {
		fmt.Println(clr.Colorize(">>Gost>> Could not find a config file. Please run gost init first.", "red"))
		return
	}
	if strings.HasSuffix(configPath, ".toml") {
		config, err := LoadFromTOML(configPath)
		if err != nil {
			fmt.Println(clr.Colorize(">>Gost>> Error loading config from .gost.toml file", "red"))
			return
		}
		fmt.Println(clr.Colorize(fmt.Sprintf("%+v", config), "green"))
	} else if strings.HasSuffix(configPath, ".json") {
		config, err := LoadFromJSON(configPath)
		if err != nil {
			fmt.Println(clr.Colorize(">>Gost>> Error loading config from .gost.json file", "red"))
			return
		}
		fmt.Println(clr.Colorize(fmt.Sprintf("%+v", config), "green"))
	} else {
		config, err := LoadFromEnv(configPath)
		if err != nil {
			fmt.Println(clr.Colorize(">>Gost>> Error loading config from .gost.env file", "red"))
			return
		}
		fmt.Println(clr.Colorize(fmt.Sprintf("%+v", config), "green"))
	}

}
